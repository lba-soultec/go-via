#  soulTec Deployment Appliance aka soulDeploy!

## Project Goal

- Automatisiertes Deployment von ESXi inkl Grundkonfiguration
- Packetierte Appliance (PhotonOS based, Installation beim Kunden via OVA)

Ablauf bei Kundeninstallation:

- Import OVA: Angabe von IP, GW, DNS der Appliance
- Erfassung des DHCP Scopes
- Upload ESXi Images
- Erfassen Host Group
- Erfassen Hosts mit MAC Adressen 

**End-Goal:**

Anhand MAC Adressen der iLOs und Values Files (ESXi Management Subnetz) werden Host Gruppen und Hosts im GoVIA automatisch erfasst.

## Tech

Die Appliance basiert auf PhotonOS und beinhaltet (goVIA)[https://github.com/maxiepax/go-via] als docker-container

goVIA stellt ein WebUI zur Verfügung. Nebst dem ist der goVIA Container dhcp und tfpd server.


## Architecture

![Architecture](https://gitlab.soultec.ch/soultec/souldeploy/-/raw/main/architecture/govia-overview.png)

-> INFOS für Netzwerk-Team des Kunden:



## Open Tasks

Please see: https://linear.app/soultec/project/soultec-deployment-appliance-a273c6e54bec/overview


### DEMO ENV

goVIA ist als docker container auf einer photonOS prepacked.

- https://10.177.176.17:8443/

- https://souldeploy.soultec.lab:8443/

*Demo Hosts (w/ iLO5):*

https://10.24.12.201/

https://10.24.12.202/


## Container Image

Das neuste Container Image ist auf diesem Repo verfügbar.
- https://gitlab.soultec.ch/soultec/souldeploy/container_registry

Die Container Registry ist eingeschränkt im Internet verfügbar. Zusätzliches Whitelisting von Public IPs können bei der internen IT beantragt werden.

## Notes on Network Config and HPE Server.

- The Server must boot from the NIC where ESXi Management Uplink should be 
- The MAC Adress of this pNIC must be added in the soulDeploy frontend
- If the ESXi Management VLAN is only tagged on the physical Network, you can set VLAN ID in the RBSU, [see here](https://support.hpe.com/hpesc/public/docDisplay?docId=a00112581en_usen_us&page=GUID-D7147C7F-2016-0901-0A69-000000000AA1.html&docLocale=en_US)
    - "Use the VLAN Configuration option to configure global VLAN settings for all enabled network interfaces. The configuration includes interfaces used in PXE boot, iSCSI boot, and HTTP/HTTPS boot, and for all preboot network access from the Embedded UEFI Shell."
    - After VLAN Config is set on RBSU, reboot the Server.


## Weitere Infos

PowerPoints und Demo Videos auf: https://soultecag.sharepoint.com/:f:/s/soulTecStreams/Eku6zA_K_BNAq2r_q4tpSeUB6BEd-oN-CDnSu2oMuw7Cbg?e=brPnlI