package rest

import (
    "fmt"
    "github.com/cosmos/cosmos-sdk/client/context"
    "github.com/gorilla/mux"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, storeName string) {
    r.HandleFunc(fmt.Sprintf("/%s/application/{%s}", storeName, "accessToken"), showApplicationsHandler(cliCtx, storeName)).Methods("GET")
    r.HandleFunc(fmt.Sprintf("/%s/application/{%s}/{%s}", storeName, "accessToken", "appCode"), showApplicationHandler(cliCtx, storeName)).Methods("GET")
    r.HandleFunc(fmt.Sprintf("/%s/admin_apps/{%s}", storeName, "accessToken"), showAdminAppsHandler(cliCtx, storeName)).Methods("GET")
    r.HandleFunc(fmt.Sprintf("/%s/is_app_user/{%s}/{%s}", storeName, "accessToken", "appCode"), showIsAppUserHandler(cliCtx, storeName)).Methods("GET")
    r.HandleFunc(fmt.Sprintf("/%s/upload/{%s}", storeName, "accessToken"), uploadFileHandler(cliCtx)).Methods("POST")
    r.HandleFunc(fmt.Sprintf("/%s/tables", storeName), createTableHandler(cliCtx)).Methods("POST")
    r.HandleFunc(fmt.Sprintf("/%s/tables/{%s}/{%s}/{%s}", storeName, "accessToken", "appCode", "tableName"), showTableHandler(cliCtx, storeName)).Methods("GET")
    r.HandleFunc(fmt.Sprintf("/%s/tables/{%s}/{%s}", storeName, "accessToken", "appCode"), showTablesHandler(cliCtx, storeName)).Methods("GET")
    r.HandleFunc(fmt.Sprintf("/%s/table-options/{%s}/{%s}/{%s}", storeName, "accessToken", "appCode", "tableName"), showTableOptionsHandler(cliCtx, storeName)).Methods("GET")
    r.HandleFunc(fmt.Sprintf("/%s/column-options/{%s}/{%s}/{%s}/{%s}", storeName, "accessToken", "appCode", "tableName", "fieldName"), showColumnOptionsHandler(cliCtx, storeName)).Methods("GET")
    r.HandleFunc(fmt.Sprintf("/%s/find/{%s}/{%s}/{%s}/{%s}", storeName, "accessToken", "appCode", "name", "id"), showRowHandler(cliCtx, storeName)).Methods("GET")
    r.HandleFunc(fmt.Sprintf("/%s/find_by/{%s}/{%s}/{%s}/{%s}/{%s}", storeName, "accessToken", "appCode", "name", "field", "value"), showIdsByHandler(cliCtx, storeName)).Methods("GET")
    r.HandleFunc(fmt.Sprintf("/%s/find_all/{%s}/{%s}/{%s}", storeName, "accessToken", "appCode", "name"), showAllIdsHandler(cliCtx, storeName)).Methods("GET")
}
