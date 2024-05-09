package macvendor

import (
	"encoding/json"
	"os"
)

type Vendor struct {
  Name string `json:"Vendor Name"`
}

var MACVendors map[string]Vendor

func lookupVendor(macAddress string) string {
  if len(macAddress) < 8 {
    return "Invalid MAC Address (less than 8 characters)"
  }
  prefix := macAddress[:8]
  vendor, ok := MACVendors[prefix]
  if ok {
    return vendor.Name
  }
  return "Not Found"
}

func LoadMACVendors(filename string) error {
  data, err := os.ReadFile(filename)
  if err != nil {
    return err
  }

  var vendors []Vendor
  if err := json.Unmarshal(data, &vendors); err != nil {
    return err
  }

  MACVendors = make(map[string]Vendor)
  for _, vendor := range vendors {
    prefix := vendor.Name[:8]  // Assuming vendor name starts with MAC prefix
    MACVendors[prefix] = vendor
  }

  return nil
}


