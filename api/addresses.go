package api

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/imdario/mergo"
	"github.com/sirupsen/logrus"
	"gitlab.soultec.ch/soultec/souldeploy/db"
	"gitlab.soultec.ch/soultec/souldeploy/models"
	"gorm.io/gorm"
)

// ListAddresses Get a list of all addresses
// @Summary Get all addresses
// @Tags addresses
// @Accept  json
// @Produce  json
// @Success 200 {array} models.Host
// @Failure 500 {object} models.APIError
// @Router /addresses [get]
func ListAddresses(c *gin.Context) {
	var items []models.Host
	if res := db.DB.Preload("Pool").Find(&items); res.Error != nil {
		Error(c, http.StatusInternalServerError, res.Error) // 500
		return
	}

	_, err := RetrieveLeases()
	if err != nil {
		Error(c, http.StatusInternalServerError, err) // 500
		return
	}
	c.JSON(http.StatusOK, items) // 200
}

func RetrieveLeases() ([]models.Host, error) {
	var leases []models.Host
	if res := db.DB.Where("expires > datetime('now', 'utc')").Find(&leases); res.Error != nil {
		if !errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, res.Error
		}
	}

	for lease := range leases {
		// Print the lease and the mac/ip address
		logrus.WithFields(logrus.Fields{
			"mac":  leases[lease].Mac,
			"ip":   leases[lease].IP,
			"pool": leases[lease].PoolID,
		}).Info("dhcp: lease")

	}
	return leases, nil
}

// GetAddress Get an existing address
// @Summary Get an existing address
// @Tags addresses
// @Accept  json
// @Produce  json
// @Param  id path int true "Address ID"
// @Success 200 {object} models.Host
// @Failure 400 {object} models.APIError
// @Failure 404 {object} models.APIError
// @Failure 500 {object} models.APIError
// @Router /addresses/{id} [get]
func GetAddress(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		Error(c, http.StatusBadRequest, err) // 400
		return
	}

	// Load the item
	var item models.Host
	if res := db.DB.Preload("Pool").First(&item, id); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			Error(c, http.StatusNotFound, fmt.Errorf("not found")) // 404
		} else {
			Error(c, http.StatusInternalServerError, res.Error) // 500
		}
		return
	}

	c.JSON(http.StatusOK, item) // 200
}

// SearchAddress Search for an address
// @Summary Search for an address
// @Tags addresses
// @Accept  json
// @Produce  json
// @Param item body models.Host true "Fields to search for"
// @Success 200 {object} models.Host
// @Failure 400 {object} models.APIError
// @Failure 404 {object} models.APIError
// @Failure 500 {object} models.APIError
// @Router /addresses/search [post]
func SearchAddress(c *gin.Context) {
	form := make(map[string]interface{})

	if err := c.ShouldBind(&form); err != nil {
		Error(c, http.StatusBadRequest, err) // 400
		return
	}

	query := db.DB

	for k, v := range form {
		query = query.Where(k, v)
	}

	// Load the item
	var item models.Host
	if res := query.Preload("Pool").First(&item); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			Error(c, http.StatusNotFound, fmt.Errorf("not found")) // 404
		} else {
			Error(c, http.StatusInternalServerError, res.Error) // 500
		}
		return
	}

	c.JSON(http.StatusOK, item) // 200
}

// CreateAddress Create a new addresses
// @Summary Create a new addresses
// @Tags addresses
// @Accept  json
// @Produce  json
// @Param item body  models.HostAddressForm true "Add ip address"
// @Success 200 {object} models.Host
// @Failure 400 {object} models.APIError
// @Failure 500 {object} models.APIError
// @Router /addresses [post]
func CreateAddress(c *gin.Context) {
	var form models.HostAddressForm

	if err := c.ShouldBind(&form); err != nil {
		Error(c, http.StatusBadRequest, err) // 400
		return
	}

	item := models.Host{HostAddressForm: form}

	// get the pool network info to verify if this ip should be added to the pool.
	var na models.Pool
	db.DB.First(&na, "id = ?", item.PoolID)

	cidr := item.IP + "/" + strconv.Itoa(na.Netmask)
	network := na.NetAddress + "/" + strconv.Itoa(na.Netmask)

	// first check if the address is even in the network.
	_, neta, _ := net.ParseCIDR(network)
	ipb, _, _ := net.ParseCIDR(cidr)
	start := net.ParseIP(na.StartAddress)
	end := net.ParseIP(na.EndAddress)
	if neta.Contains(ipb) {
		//then check if it's in the given range by the pool.
		trial := net.ParseIP(item.IP)

		if bytes.Compare(trial, start) >= 0 && bytes.Compare(trial, end) <= 0 {
			logrus.WithFields(logrus.Fields{
				"ip":    trial.String(),
				"start": start.String(),
				"end":   end.String(),
			}).Debug("ip validation successful")
		} else {
			logrus.WithFields(logrus.Fields{
				"ip":    trial.String(),
				"start": start.String(),
				"end":   end.String(),
			}).Debug("the ip address is not in the scope of the dhcp pool associated with the group")
			Error(c, http.StatusBadRequest, fmt.Errorf("the ip address is not in the scope of the dhcp pool associated with the group")) // 400
			return
		}
	} else {
		logrus.WithFields(logrus.Fields{
			"netmask": na.Netmask,
			"start":   start.String(),
			"end":     end.String(),
			"ip":      item.IP,
		}).Debug("the ip address is not in the scope of the dhcp pool associated with the group")
		Error(c, http.StatusBadRequest, fmt.Errorf("the ip address is not in the scope of the dhcp pool associated with the group")) // 400
		return
	}

	// ensure the mac address is properly formated.
	mac, _ := net.ParseMAC(item.Mac)
	item.Mac = mac.String()

	// if ip address checks pas, continue to commit.
	if item.ID != 0 { // Save if its an existing item
		if res := db.DB.Save(&item); res.Error != nil {
			Error(c, http.StatusInternalServerError, res.Error) // 500
			return
		}
	} else { // Create a new item
		if res := db.DB.Create(&item); res.Error != nil {
			Error(c, http.StatusInternalServerError, res.Error) // 500
			return
		}
	}

	// Load a new version with relations
	if res := db.DB.Preload("Pool").First(&item); res.Error != nil {
		Error(c, http.StatusInternalServerError, res.Error) // 500
		return
	}

	c.JSON(http.StatusOK, item) // 200

	logrus.WithFields(logrus.Fields{
		"Hostname": item.Hostname,
		"Domain":   item.Domain,
		"IP":       item.IP,
		"MAC":      item.Mac,
		"Pool ID":  item.PoolID,
		"Group ID": item.GroupID,
	}).Debug("host")
}

// UpdateAddress Update an existing address
// @Summary Update an existing address
// @Tags addresses
// @Accept  json
// @Produce  json
// @Param  id path int true "Address ID"
// @Param  item body  models.HostAddressForm true "Update an ip address"
// @Success 200 {object} models.Host
// @Failure 400 {object} models.APIError
// @Failure 404 {object} models.APIError
// @Failure 500 {object} models.APIError
// @Router /addresses/{id} [patch]
func UpdateAddress(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		Error(c, http.StatusBadRequest, err) // 400
		return
	}

	// Load the form data
	var form models.HostAddressForm
	if err := c.ShouldBind(&form); err != nil {
		Error(c, http.StatusBadRequest, err) // 400
		return
	}

	// Load the item
	var item models.Host
	if res := db.DB.First(&item, id); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			Error(c, http.StatusNotFound, fmt.Errorf("not found")) // 404
		} else {
			Error(c, http.StatusInternalServerError, res.Error) // 500
		}
		return
	}

	// Merge the item and the form data
	if err := mergo.Merge(&item, models.Host{HostAddressForm: form}, mergo.WithOverride); err != nil {
		Error(c, http.StatusInternalServerError, err) // 500
	}

	// Mergo doesn't overwrite 0 or false values, force set
	item.Reimage = form.Reimage
	item.Progress = form.Progress

	// Save it
	if res := db.DB.Save(&item); res.Error != nil {
		Error(c, http.StatusInternalServerError, res.Error) // 500
		return
	}

	// Load a new version with relations
	if res := db.DB.Preload("Pool").First(&item); res.Error != nil {
		Error(c, http.StatusInternalServerError, res.Error) // 500
		return
	}

	c.JSON(http.StatusOK, item) // 200
}

// DeleteAddress Remove an existing address
// @Summary Remove an existing address
// @Tags addresses
// @Accept  json
// @Produce  json
// @Param  id path int true "Address ID"
// @Success 204
// @Failure 404 {object} models.APIError
// @Failure 500 {object} models.APIError
// @Router /addresses/{id} [delete]
func DeleteAddress(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		Error(c, http.StatusBadRequest, err) // 400
		return
	}

	// Load the item
	var item models.Host
	if res := db.DB.First(&item, id); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			Error(c, http.StatusNotFound, fmt.Errorf("not found")) // 404
		} else {
			Error(c, http.StatusInternalServerError, res.Error) // 500
		}
		return
	}

	// delete it
	if res := db.DB.Delete(&item); res.Error != nil {
		Error(c, http.StatusInternalServerError, res.Error) // 500
		return
	}

	c.JSON(http.StatusNoContent, gin.H{}) //204
}
