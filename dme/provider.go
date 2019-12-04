package dme

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"os"
)

// Provider provides a Provider...
func Provider() terraform.ResourceProvider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"akey": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: envDefaultFunc("DME_AKEY"),
				Description: "A DNSMadeEasy API Key.",
			},
			"skey": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: envDefaultFunc("DME_SKEY"),
				Description: "The Secret Key for API operations.",
			},
			"usesandbox": &schema.Schema{
				Type:        schema.TypeBool,
				Required:    true,
				DefaultFunc: envDefaultFunc("DME_USESANDBOX"),
				Description: "If true, use the DME Sandbox.",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"dme_record": resourceDMERecord(),
		},
	}

	provider.ConfigureFunc = func(d *schema.ResourceData) (interface{}, error) {
		terraformVersion := provider.TerraformVersion
		if terraformVersion == "" {
			// Terraform 0.12 introduced this field to the protocol
			// We can therefore assume that if it's missing it's 0.10 or 0.11
			terraformVersion = "0.11+compatible"
		}
		return providerConfigure(d, terraformVersion)
	}

	return provider

}

func envDefaultFunc(k string) schema.SchemaDefaultFunc {
	return func() (interface{}, error) {
		if v := os.Getenv(k); v != "" {
			if v == "true" {
				return true, nil
			} else if v == "false" {
				return false, nil
			}
			return v, nil
		}
		return nil, nil
	}
}

func providerConfigure(d *schema.ResourceData, terraformVersion string) (interface{}, error) {
	config := Config{
		AKey:       d.Get("akey").(string),
		SKey:       d.Get("skey").(string),
		UseSandbox: d.Get("usesandbox").(bool),
	}
	return config.Client()
}
