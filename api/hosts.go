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
	"github.com/maxiepax/go-via/db"
	"github.com/maxiepax/go-via/models"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// ListHosts Get a list of all hosts
// @Summary Get all hosts
// @Tags hosts
// @Accept  json
// @Produce  json
// @Success 200 {array} models.Address
// @Failure 500 {object} models.APIError
// @Router /hosts [get]
func ListHosts(c *gin.Context) {
	var items []models.Host
	if res := db.DB.Preload("Pool").Find(&items); res.Error != nil {
		Error(c, http.StatusInternalServerError, res.Error) // 500
		return
	}
	c.JSON(http.StatusOK, items) // 200
}

// GetHost Get an existing Host
// @Summary Get an existing Host
// @Tags hosts
// @Accept  json
// @Produce  json
// @Param  id path int true "Host ID"
// @Success 200 {object} models.Host
// @Failure 400 {object} models.APIError
// @Failure 404 {object} models.APIError
// @Failure 500 {object} models.APIError
// @Router /hosts/{id} [get]
func GetHost(c *gin.Context) {
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

// SearchHost Search for an host
// @Summary Search for an host
// @Tags hosts
// @Accept  json
// @Produce  json
// @Param item body models.Host true "Fields to search for"
// @Success 200 {object} models.Host
// @Failure 400 {object} models.APIError
// @Failure 404 {object} models.APIError
// @Failure 500 {object} models.APIError
// @Router /hosts/search [post]
func SearchHost(c *gin.Context) {
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

// CreateHost Create a new host
// @Summary Create a new host
// @Tags hosts
// @Accept  json
// @Produce  json
// @Param item body models.HostForm true "Add a host"
// @Success 200 {object} models.Host
// @Failure 400 {object} models.APIError
// @Failure 500 {object} models.APIError
// @Router /hosts [post]
func CreateHost(c *gin.Context) {
	var form models.HostForm

	if err := c.ShouldBind(&form); err != nil {
		Error(c, http.StatusBadRequest, err) // 400
		return
	}

	item := models.Host{HostForm: form}

	// get the pool network info to verify if this ip should be added to the pool.
	var na models.Pool
	db.DB.First(&na, "id = ?", item.HostForm.PoolID)

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
				"ip":    trial,
				"start": start,
				"end":   end,
			}).Debug("ip validation successful")
		} else {
			logrus.WithFields(logrus.Fields{
				"ip":    trial,
				"start": start,
				"end":   end,
			}).Debug("the ip address is not in the scope of the dhcp pool associated with the group")
			Error(c, http.StatusBadRequest, fmt.Errorf("the ip address is not in the scope of the dhcp pool associated with the group")) // 400
			return
		}
	} else {
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

// UpdateHost Update an existing host
// @Summary Update an existing host
// @Tags hosts
// @Accept  json
// @Produce  json
// @Param  id path int true "Host ID"
// @Param  item body models.HostForm true "Update a host"
// @Success 200 {object} models.Host
// @Failure 400 {object} models.APIError
// @Failure 404 {object} models.APIError
// @Failure 500 {object} models.APIError
// @Router /hosts/{id} [patch]
func UpdateHost(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		Error(c, http.StatusBadRequest, err) // 400
		return
	}

	// Load the form data
	var form models.HostForm
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
	if err := mergo.Merge(&item, models.Host{HostForm: form}, mergo.WithOverride); err != nil {
		Error(c, http.StatusInternalServerError, err) // 500
	}

	// Mergo doesn't overwrite 0 or false values, force set
	item.HostForm.Reimage = form.Reimage
	item.HostForm.Progress = form.Progress

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

// DeleteHost Remove an existing host
// @Summary Remove an existing host
// @Tags hosts
// @Accept  json
// @Produce  json
// @Param  id path int true "Host ID"
// @Success 204
// @Failure 404 {object} models.APIError
// @Failure 500 {object} models.APIError
// @Router /hosts/{id} [delete]
func DeleteHost(c *gin.Context) {
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
