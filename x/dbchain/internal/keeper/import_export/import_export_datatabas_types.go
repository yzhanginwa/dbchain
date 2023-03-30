package import_export 

type Database struct {
    Name string    `json:"name"`
    Appcode string `json:"appCode"`
    Memo string    `json:"memo"`
    Tables []Table `json:"table"`

    Kind int                  `json:"type"`
    PermissionReq bool        `json:"permission_required"`
    CustomFns []CustomFn      `json:"custom_fns"`
    CustomQueriers []CustomFn `json:"custom_queriers"`
}

type Table struct {
    Name string      `json:"name"`
    Memo string      `json:"memo"`
    Fields []Field   `json:"field"`
    Filter string    `json:"filter"`
    Trigger string   `json:"trigger"`
    Options []string `json:"options"`
}

type Field struct {
    Name string          `json:"name"`
    Memo string          `json:"memo"`
    FieldType string     `json:"fieldType"`
    PropertyArr []string `json:"propertyArr"`
    IsIndex bool         `json:"inIndex"`
}

type CustomFn struct {
    Name string       `json:"name"`
    Owner string      `json:"owner"`
    Description string `json:"description"`
    Body string       `json:"body"`
}

