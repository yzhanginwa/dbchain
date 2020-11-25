package super_script

/////////////////////////////////////////////////////////////////////////////////
//                                                                             //
//                                     BNF                                     //
//                                                                             //
/////////////////////////////////////////////////////////////////////////////////
//                                                                             //
//    statements = statement [ statements]                                     //
//                                                                             //
//    statement = if_condition | return | insert .                             //
//                                                                             //
//    return = "return" "(" ( "true" | "false" ) ")" .                         //
//                                                                             //
//    insert = "insert" "(" table_name "," field_name ", " single_value        //
//             [ "," field_name "," single_value ] ")" .                       //
//                                                                             //
//    if_condition = "if" "(" condition ")" "{" statements "}"                 //
//                   [ "else" "{" statements "}" ] .                           //
//                                                                             //
//    condition = exist | comparison .                                         //
//                                                                             //
//    comparison = single_value "==" single_value                              //
//               | single_value "in" list_value .                              //
//                                                                             //
//    exist = "exist" "(" table_value ")"                                      //
//                                                                             //
//    single_value = this_expr                                                 //
//                 | string_literal .                                          //
//                                                                             //
//    this_expr = "this" "." field [ "." parent_field ]                        //
//                                                                             //
//    string_literal = double_quote ident double_quote .                       //
//                                                                             //
//    parent_field = "parent" "." field { "." parent_field } .                 //
//                                                                             //
//    list_value = "(" string_literal ["," string_listeral] ")" .              //
//                                                                             //
//    table_value = "table" "." table_name { "." where } .                     //
//                                                                             //
//    table_name = ident .                                                     //
//                                                                             //
//    where = "where" "(" field "==" single_value ")" .                        //
//                                                                             //
//    field = ident .                                                          //
//                                                                             //
//    ident = (a-zA-Z_) .                                                      //
//                                                                             //
/////////////////////////////////////////////////////////////////////////////////
