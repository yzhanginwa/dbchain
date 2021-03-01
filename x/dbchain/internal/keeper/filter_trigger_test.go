package keeper

import (
	lua "github.com/yuin/gopher-lua"
	ss "github.com/yzhanginwa/dbchain/x/dbchain/internal/super_script"
	"strings"
	"testing"
)

func TestFilterTrigger(t *testing.T) {

	L := lua.NewState()
	defer L.Close()


	L.SetGlobal("Insert", L.NewFunction(Insert))
	L.SetGlobal("fieldIn", L.NewFunction(fieldIn))
	L.SetGlobal("exist", L.NewFunction(exist))

	owner := "cosmos1a5jps5hyu2n0zjdzlkh8llz9efp9g75f45zmqr"
	filters := []string {
		`if (this.corp_id.parent.created_by == this.created_by){ 
			return(true) 
		}`,

		`if (this.name in ("Bob", "bar")){ 
			return(true) 
		} else {
			return(false)
		}`,

		`if (exist(table.corp.where(id == this.corp_id))){ 
			return(true) 
		}`,

		`if (exist(table.corp.where(type == "corp").where(name == "Microsoft"))) { 
			return(true) 
		}`,

		`if (this.name in ("foo", "bar")) { 
			return(false) 
		} elseif (this.name in ("Bob", "bar")) {  
			if(this.name in ("Bob")){ 
				return(true)
			} 
		}else { 
			return (false) 
		}`,
		`if (this.name in ("foo", "bar")) {  
           if (this.name in ("foo", "bar")) { 
				return (false) 
			} elseif (this.name in ("foo", "bar")) { 
				return (false) 
			}
         }
        if (this.name in ("foo", "bar")) {  
          	if (this.name in ("foo", "bar")) { 
				return (false) 
			}elseif (this.name in ("foo", "bar")) { 
				return (false) 
			}
        } elseif (this.name in ("Bob", "bar")) {  
			if (this.name in ("foo", "bar")) { 
				return(false) 
			}elseif (this.name in ("Bob", "bar")) { 
				if (this.name in ("foo", "bar")) { 
					return (false) 
				}elseif (this.name in ("Bob", "bar")) { 
					return (true) 
				}else {
					return (false) 
				}
          }
        } else { 
			return (true) 
		}`,


		//trigger
		`if (this.corp_id.parent.created_by == this.created_by){ 
			ret = Insert("tabname","val") 
			return ret
		}`,
	}
	//
	result := []string {
		"true","true","true","true","true","true","true",
	}
	for index, val := range filters {
		registerThisData(owner,L, val)
		p := ss.NewPreprocessor(strings.NewReader(val))
		p.Process()
		newScript := p.Reconstruct()
		err := L.DoString(newScript)
		if err != nil {
			t.Errorf("index : %v err : %v", index,err)
			continue
		}
		ret := L.Get(1).(lua.LBool)
		if ret.String() != result[index]{
			t.Errorf("want : %v, actual %v\n",result[index],ret.String())
			t.Errorf("%d failed", index)
		}
	}

}



func registerThisData(owner string, L *lua.LState,script string) bool {
	this := L.NewTable()
	tbFields := []string {
		"id",
		"create_at",
		"create_by",
		"corp_id",
		"name",
		"age",
		"Tel",
	}

	fieldsVal := map[string]string {
		"id" : "1",
		"create_at" : "2021-02-26|17:16:39.986",
		"create_by" : "cosmos1a5jps5hyu2n0zjdzlkh8llz9efp9g75f45zmqr",
		"corp_id" : "1",
		"name" : "Bob",
		"age" : "19",
		"Tel" : "13112345678",
	}

	parentFieldVal := map[string]string {
		"id" : "1",
		"create_at" : "2021-02-26|17:16:39.986",
		"create_by" : "cosmos1a5jps5hyu2n0zjdzlkh8llz9efp9g75f45zmqr",
		"occupation" : "worker",
	}

	//register func need to be changed like this
	for _, field := range tbFields {
		isForeignKey := false
		v , ok := fieldsVal[field]
		if ok { //is foreignKey
				//if !strings.HasSuffix(field, "_id"){
				//this.RawSetString(field, lua.LString(v))
				if strings.Contains(field,"_") {
					temp := "." + field + ".parent"
					if strings.Contains(script,temp) {
						isForeignKey = true
					}
				}
				if !isForeignKey {
					this.RawSetString(field, lua.LString(v))
				}
		} else if field == "created_by" {
			this.RawSetString(field, lua.LString(owner))
		} else {
			this.RawSetString(field, lua.LString("nil"))
		}

		if isForeignKey {// has parent table and the value is not null
			foreignKey := L.NewTable()
			parent := L.NewTable()
			//parentTableName := field[:len(field)-3]
			//tableAndKey := strings.Split(field,"_")
			//if len(tableAndKey) != 2 {
			//	return false
			//}
			//
			//parentTableName := tableAndKey[0]
			//parentTableField := tableAndKey[1]
			//
			//ids := k.FindBy(ctx, appId, parentTableName, parentTableField, []string{v}, owner)
			//if len(ids) != 1 {
			//	return false
			//}
			//RowFields, err := k.DoFind(ctx, appId, parentTableName, ids[0])
			//if err != nil {
			//	return false
			//}
			//for key , value  := range RowFields {
			//	parent.RawSetString(key, lua.LString(value))
			//}

			for key ,val := range parentFieldVal {
				parent.RawSetString(key, lua.LString(val))
			}
			foreignKey.RawSetString("parent", parent)
			this.RawSetString(field, foreignKey)
		}
	}
	L.SetGlobal("this", this) //register this table
	return true
}
//this fun only is used to test
func Insert(L *lua.LState) int {
	L.Push(lua.LTrue)
	return 1
}

//this fun only is used to test
func fieldIn(L *lua.LState) int {
	ParamsNum := L.GetTop()
	if ParamsNum < 2 {
		L.Push(lua.LBool(false))
		return 1
	}

	src := L.ToString(1)
	for i := 2; i <= ParamsNum; i++ {
		dst := L.ToString(i)
		if src == dst{
			L.Push(lua.LBool(true))
			return 1
		}
	}
	L.Push(lua.LBool(false))
	return 1
}
//this fun only is used to test
func exist(L *lua.LState) int {
	L.Push(lua.LTrue)
	return 1
}