package paging

import (
	"strings"
	"strconv"
	"fmt"

	"github.com/yuanchi/paging/util"
)

/*
referring to 'https://www.slideshare.net/slideshow/view?login=Eweaver&preview=no&slideid=1&title=efficient-pagination-using-mysql'
to implement the function of pagination
*/
type Paging struct {
	Id string
	IdVal interface{}
	IdDesc bool

	Next bool
	Limit int
	Fields []*FieldData
}

type FieldData struct {
        Name string
        Value interface{}
        Desc, Unique bool
}

type SortDirect interface {
	Next() string
	Prev() string
}

const (
	gt = ">"
	lt = "<"
	ge = ">="
	le = "<="
)

type uniDesc struct {}
type uniAsc struct {}
type dupDesc struct {}
type dupAsc struct {}

func (ud *uniDesc) Next() string { return lt }
func (ud *uniDesc) Prev() string { return gt }

func (ua *uniAsc) Next() string { return gt }
func (ua *uniAsc) Prev() string { return lt}

func (dd *dupDesc) Next() string { return le }
func (dd *dupDesc) Prev() string { return ge}

func (da *dupAsc) Next() string { return ge }
func (da *dupAsc) Prev() string { return le }

var (
	ud = uniDesc{}
	ua = uniAsc{}
	dd = dupDesc{}
	da = dupAsc{}
)
 
func ToSortDirect(unique, desc bool) SortDirect {
	if unique && desc {
		return &ud
	} else if unique && !desc {
		return &ua
	} else if !unique && desc {
		return &dd
	} else {
		return &da
	}
}

// nextTpl needs three parameters: {{.Conds}}, {{.Sorts}}, {{.Limit}}
func SortQueryByUniqueKey(p *Paging, nextTpl, alias string) string {
	id := p.IdVal
	var idCond string
	idName := alias + "." + p.Id + " "
	if id == nil {
		idCond = "1 = 1"
	} else {
		sort := ToSortDirect(true, p.IdDesc)
        	var direct string
        	if p.Next {
                	direct = sort.Next()
        	} else {
                	direct = sort.Prev()
       		 }
		idCond = idName + direct + " " + PrepareToString(id)
	}
	
	desc := "DESC"
        if !p.IdDesc {
                desc = "ASC"
        }
	idSort := idName + desc

	sql := strings.Replace(nextTpl, "{{.Limit}}", strconv.Itoa(p.Limit), -1)
	sql = strings.Replace(sql, "{{.Conds}}", idCond, -1)
	sql = strings.Replace(sql, "{{.Sorts}}", idSort, -1)	

	return sql
}

func PrepareToString(val interface{}) string {
	var s string
	switch v := val.(type) {
		case string:
			s = util.PrepareString(v)
		case int:
			s = strconv.Itoa(v)
		default:
			s = fmt.Sprintf("%+v", v)
			fmt.Printf("PrepareToString val is %+v", v)
	}
	return s
}

func SortQueryBy(p *Paging, nextTpl, alias string) string {
	id := p.IdVal
	var idDesc string
	if p.IdDesc {
		idDesc = "DESC"
	} else {
		idDesc = "ASC"
	}
	next := p.Next	

	idName := alias + "." + p.Id
	idSortDirect := ToSortDirect(true, p.IdDesc)
	var idDirect string
	if next {
		idDirect = idSortDirect.Next()
	} else {
		idDirect = idSortDirect.Prev()
	}

	var idCond string 
	if id != nil {
		idCond = idName + " " + idDirect + " " + PrepareToString(id)
	} else {
		idCond = "1 = 1"
	}
	idSort := idName + " " + idDesc

	var where []string
	var orderBy []string
	var containsDuplicable bool
	for _, f := range p.Fields {
		if f.Unique {
                        continue
                }

		sort := ToSortDirect(f.Unique, f.Desc)
		condName := alias + "." + f.Name + " "

		sortItem := condName
		if f.Desc {
			sortItem += "DESC"
		} else {
			sortItem += "ASC"
		}		
		orderBy = append(orderBy, sortItem)				
		
		if f.Value == nil { // if there's no value, representing need extra fields to sort, but not initialized yet  
			continue
		}
		containsDuplicable = true
		cond := condName	
		if next {
			cond += sort.Next()
		} else {
			cond += sort.Prev()
		}
		cond += " "
		cond += PrepareToString(f.Value)
		where = append(where, cond)
		
	}
	orderBy = append(orderBy, idSort)

	extraCond := idCond
	var conds string
	for _, cond := range where {
		extraCond += " OR "
		if strings.Contains(cond, ge) {
			extraCond += strings.Replace(cond, ge, gt, -1)
		} else {
			extraCond += strings.Replace(cond, le, lt, -1)
		}
		conds += (cond + " AND ")
	}
	if containsDuplicable {
		conds += ("(" + extraCond + ")")
	} else {
		conds = idCond
	}
	sql := strings.Replace(nextTpl, "{{.Conds}}", conds, -1)

	sorts := strings.Join(orderBy, ", ")
	sql = strings.Replace(sql, "{{.Sorts}}", sorts, -1)

	sql = strings.Replace(sql, "{{.Limit}}", strconv.Itoa(p.Limit), -1)
	return sql
}