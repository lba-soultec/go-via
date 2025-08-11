package api

import (
	"errors"
	"fmt"
	"strconv"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gitlab.soultec.ch/soultec/souldeploy/ilomapi"
)

func StartHost(c *gin.Context) {
	// Log that the method has been called
	logrus.Debug("Received start request")
	// Extract query parameters and create API client
	api, _, err := createAPIClientFromParams(c)
	if err != nil {
		return
	}

	err = api.StartServer()

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"endpoint": api.GetEndpoint(),
			"error":    err.Error(),
		}).Error("Failed to start host")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to start host",
			"message": err.Error(),
		})
		return
	}
	// Log the host configuration
	logrus.WithFields(logrus.Fields{
		"endpoint": api.GetEndpoint(),
	}).Info("Host has been scheduled to start successfully")
	// Return the host configuration as JSON
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Host configuration fetched successfully",

		"endpoint":   api.GetEndpoint(),
		"apiFlavour": api.GetFlavour(),
	})

}

func createAPIClientFromParams(c *gin.Context) (ilomapi.IlomApi, map[string]string, error) {
	// Extract query parameters
	var params map[string]string
	var requestBody struct {
		IloIpAddr  string `json:"iloIpAddr"`
		Port       string `json:"port"`
		ApiFlavour string `json:"apiFlavour"`
		Username   string `json:"username"`
		Password   string `json:"password"`
		VlanID     int    `json:"vlanID"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {

		fmt.Println(c.Request.PostForm)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"message": err.Error(),
		})
		return nil, params, errors.New("Invalid request body: " + err.Error())
	}

	iloIpAddr := requestBody.IloIpAddr
	port := requestBody.Port
	apiFlavour := requestBody.ApiFlavour
	username := requestBody.Username
	password := requestBody.Password

	params = map[string]string{
		"iloIpAddr":  requestBody.IloIpAddr,
		"port":       requestBody.Port,
		"apiFlavour": requestBody.ApiFlavour,
		"username":   requestBody.Username,
		"password":   requestBody.Password,
		"vlanID":     fmt.Sprintf("%d", requestBody.VlanID),
	}

	var api ilomapi.IlomApi

	// Validate parameters
	if iloIpAddr == "" || port == "" || apiFlavour == "" || username == "" || password == "" {
		if len(password) > 0 {
			password = "*****" // Mask password for security
		}
		logrus.WithFields(logrus.Fields{
			"iloIpAddr":  iloIpAddr,
			"port":       port,
			"apiFlavour": apiFlavour,
			"username":   username,
			"password":   password, // Mask password for security
		}).Warn("Missing required parameters")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "iloIpAddr, apiFlavour, port, username, and password are required parameters",
		})
		return nil, params, errors.New("Missing required parameters: iloIpAddr:" + iloIpAddr +
			" port:" + port +
			" apiFlavour:" + apiFlavour +
			" username:" + username +
			" password:" + password)
	}

	// Get the IP and mac for every Network Interface

	switch apiFlavour {

	case "redfish":
		// Call the Redfish API to get the host configuration
		api = ilomapi.NewRedFishApi(iloIpAddr, port, username, password)
	}
	// return the API client
	return api, params, nil
}

func ShutdownHost(c *gin.Context) {

	// Log that the method has been called
	logrus.Debug("Received shutdown request")

	// Extract query parameters and create API client
	api, _, err := createAPIClientFromParams(c)
	if err != nil {
		return
	}

	err = api.StopServer()

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"endpoint": api.GetEndpoint(),
			"error":    err.Error(),
		}).Error("Failed to shutdown host")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to shutdown host",
			"message": err.Error(),
		})
		return
	}
	// Log the host configuration
	logrus.WithFields(logrus.Fields{
		"endpoint": api.GetEndpoint(),
	}).Debug("Host has been shut down successfully")
	// Return the host configuration as JSON
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Host configuration fetched successfully",

		"endpoint":   api.GetEndpoint(),
		"apiFlavour": api.GetFlavour(),
	})

}

func RebootHost(c *gin.Context) {
	// Log that the method has been called
	logrus.Debug("Received reboot request")
	// Extract query parameters and create API client
	api, _, err := createAPIClientFromParams(c)
	if err != nil {
		return
	}

	err = api.RebootServer()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"endpoint": api.GetEndpoint(),
			"error":    err.Error(),
		}).Error("Failed to reboot host")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to reboot host",
			"message": err.Error(),
		})
		return
	}

	// Log success and return response
	logrus.WithFields(logrus.Fields{
		"endpoint": api.GetEndpoint(),
	}).Info("Host has been rebooted successfully")
	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"message":    "Host rebooted successfully",
		"endpoint":   api.GetEndpoint(),
		"apiFlavour": api.GetFlavour(),
	})
}

func OneTimeBoot(c *gin.Context) {

	// Log that the method has been called
	logrus.Info("Received reboot request")
	// Extract query parameters and create API client
	api, _, err := createAPIClientFromParams(c)
	if err != nil {
		return
	}

	err = api.SetOneTimeHTTPBoot()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"endpoint": api.GetEndpoint(),
			"error":    err.Error(),
		}).Error("Failed to reboot host")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to reboot host",
			"message": err.Error(),
		})
		return
	}

	// Log success and return response
	logrus.WithFields(logrus.Fields{
		"endpoint": api.GetEndpoint(),
	}).Info("Host has been rebooted successfully")
	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"message":    "Host rebooted successfully",
		"endpoint":   api.GetEndpoint(),
		"apiFlavour": api.GetFlavour(),
	})
}

func SetVLANID(c *gin.Context) {
	// Log that the method has been called
	logrus.Debug("Received set VLANID request")

	// Extract query parameters and create API client
	api, params, err := createAPIClientFromParams(c)
	if err != nil {
		return
	}

	vlanID, err := strconv.Atoi(params["vlanID"])

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"endpoint": api.GetEndpoint(),
			"error":    err.Error(),
		}).Error("Failed to parse VLAN ID")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to parse VLAN ID",
			"message": err.Error(),
		})
		return
	}
	// Set VLAN ID
	err = api.SetVLANID(vlanID)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"endpoint": api.GetEndpoint(),
			"error":    err.Error(),
		}).Error("Failed to reboot host")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to reboot host",
			"message": err.Error(),
		})
		return
	}

	err = api.SetOneTimeHTTPBoot()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"endpoint": api.GetEndpoint(),
			"error":    err.Error(),
		}).Error("Failed to set one Time Boot host")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to set one Time Boot host",
			"message": err.Error(),
		})
		return
	}

	err = api.RebootServer()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"endpoint": api.GetEndpoint(),
			"error":    err.Error(),
		}).Error("Failed to trigger reboot host")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to trigger reboot host",
			"message": err.Error(),
		})
		return
	}

	// Log success and return response
	logrus.WithFields(logrus.Fields{
		"endpoint": api.GetEndpoint(),
	}).Info("Host has been configured with VLANID " + params["vlanID"] + " successfully")
	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"message":    "Host has been configured with VLANID " + params["vlanID"] + " successfully",
		"endpoint":   api.GetEndpoint(),
		"apiFlavour": api.GetFlavour(),
	})
}
