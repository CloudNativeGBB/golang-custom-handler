package settings

import (
	"fmt"
	"identity/structs"
)

// Tenant variable
var Tenant = "[YOUR_TENANT_NAME].onmicrosoft.com"

// TenantB2CLogin variable
var TenantB2CLogin = "[YOUR_TENANT_NAME].b2clogin.com"

// PrimaryTrustFrameworkPolicy variable
var PrimaryTrustFrameworkPolicy = "[YOUR_ROPC_FLOW_NAME]"

// ClientID variable
var ClientID = "[YOUR_APP_ID]"

// ClientSecret variable
var ClientSecret = "[YOUR_APP_ID_SECRET]"

// GlobalAuthToken variable
var GlobalAuthToken = ""

// GlobalPublicKeys variable
var GlobalPublicKeys = structs.Keys{}

func init() {
	fmt.Println("package: settings - initialized")
}
