package keeper

///////////////////////////////
//                           //
//     insert row            //
//                           //
///////////////////////////////

type ScriptInsertRow struct {
	TableName string     `json:"table_name"`
	Fields map[string]string `json:"fields"`
}

///////////////////////////////
//                           //
//   foreign insert row      //
//                           //
///////////////////////////////

type ScriptForeignInsertRow struct {
	TableName string     `json:"table_name"`
	Fields map[string]string `json:"fields"`
	ForeignKey []string     `json:"foreign_key"`
}

///////////////////////////////
//                           //
//   freeze multiple row     //
//                           //
///////////////////////////////

type ScriptFreezeMultRow struct {
	TableName string     `json:"table_name"`
	Ids       []string    `json:"ids"`
}

///////////////////////////////////////
//                                   //
//   freeze multiple row  by field   //
//                                   //
///////////////////////////////////////

type ScriptFreezeMultRowByField struct {
	ScriptInsertRow
}