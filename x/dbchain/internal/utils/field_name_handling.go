package utils

func GetTableNameFromForeignKey(fk string) (string, bool) { // fk is like supplier_id, address_id
    l := len(fk)
    if l < 4 {
        return "", false
    }

    part1 := fk[0:l-3]
    part2 := fk[l-3:l]

    if part2 == "_id" {
        return part1, true
    } else {
        return "", false
    }
}
