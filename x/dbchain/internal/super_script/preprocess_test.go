package super_script

import (
	"strings"
	"testing"
)

func TestPreprocessor(t *testing.T) {
	src := []string{
		"if (this.corp_id.parent.created_id == this.created_id) { return(true) }",
		`if (this.name in ("foo", "bar")) { return(true) }`,
		`if (exist(table.corp.where(id == this.corp_id))) { return(true) }`,
		`if (exist(table.corp.where(type == "corp").where(name == "Microsoft"))) { return(true) }`,
		`if (this.name in ("foo", "bar")) { return(true) } 
		elseif (this.name in ("foo", "bar")) {  if(this.name in ("foo", "bar")){ return(true)} }
		else { return (true) }`,
		`if (this.name in ("foo", "bar")) {  
           if (this.name in ("foo", "bar")) { return (true) }
           elseif (this.name in ("foo", "bar")) { return (true) }
         }
        if (this.name in ("foo", "bar")) {  
          if (this.name in ("foo", "bar")) { return (true) }
          elseif (this.name in ("foo", "bar")) { return (true) }
        }
        elseif (this.name in ("foo", "bar")) {  
          if (this.name in ("foo", "bar")) { return(true) } 
          elseif (this.name in ("foo", "bar")) { 
            if (this.name in ("foo", "bar")) { return (true) }
            elseif (this.name in ("foo", "bar")) { return (true) }
            else { return (true) }
          }
        }
        else { return (true) }`,
	}

	dst := []string{
		`if(this.corp_id.parent.created_id==this.created_id) then  return(true)  end `,
		`if(fieldIn(this.name,"foo","bar")) then  return(true)  end `,
		`if(exist("corp","id","==",this.corp_id)) then  return(true)  end `,
		`if(exist("corp","type","==","corp","name","==","Microsoft")) then  return(true)  end `,
		`if(fieldIn(this.name,"foo","bar")) then  return(true)  
		elseif(fieldIn(this.name,"foo","bar")) then   if(fieldIn(this.name,"foo","bar")) then  return(true) end  
		else return (true)  end `,
		`if(fieldIn(this.name,"foo","bar")) then   
           if(fieldIn(this.name,"foo","bar")) then  return (true) 
           elseif(fieldIn(this.name,"foo","bar")) then  return (true)  end 
          end 
        if(fieldIn(this.name,"foo","bar")) then   
          if(fieldIn(this.name,"foo","bar")) then  return (true) 
          elseif(fieldIn(this.name,"foo","bar")) then  return (true)  end 
        
        elseif(fieldIn(this.name,"foo","bar")) then   
          if(fieldIn(this.name,"foo","bar")) then  return(true)  
          elseif(fieldIn(this.name,"foo","bar")) then  
            if(fieldIn(this.name,"foo","bar")) then  return (true) 
            elseif(fieldIn(this.name,"foo","bar")) then  return (true) 
            else return (true)  end 
           end 
        
        else return (true)  end `,
	}
	for i,v := range src{
		p := NewPreprocessorOld(strings.NewReader(v))
		p.Process()
		s := p.Reconstruct()
		if s != dst[i] {
			t.Errorf("%d failed", i)
		}
	}
}
