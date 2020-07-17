package super_script

/////////////////////////////////////////////////////////////////////////////////
//                                                                             //
//                                     BNF                                     //
//                                                                             //
/////////////////////////////////////////////////////////////////////////////////
//                                                                             //
//    statement = if_condition | return | insert .                             //
//                                                                             //
//    return = "return" "(" ( "true" | "false" ) ")" .                         //
//                                                                             //
//    insert = "insert" "(" table_name "," field_name ", " string_literal      //
//             [ "," field_name "," string_literal ] ")" .                     //
//                                                                             //
//    if_condition = "if" filter_condition "then" statement "fi" .             //
//                                                                             //
//    filter_condition = single_value "==" single_value .                      //
//                     | single_value "in" multi_value .                       //
//                                                                             //
//    single_value = this_expr                                                 //
//                 | string_literal .                                          //
//                                                                             //
//    this_expr = this "." field [ "." parent_field ]                          //
//                                                                             //
//    string_literal = double_quote ident double_quote .                       //
//                                                                             //
//    parent_field = parent "." field { "." parent_field } .                   //
//                                                                             //
//    multi_value = list_value | table_value .                                 //
//                                                                             //
//    list_value = "(" string_literal ["," string_listeral] ")" .              //
//                                                                             //
//    table_value = table "." table_name { "." where } "." field               //
//                | "(" string_literal { "," string_literal } ")" .            //
//                                                                             //
//    table_name = ident .                                                     //
//                                                                             //
//    where = "where" "(" field "==" single_value ")" .                        //
//                                                                             //
//    field = ident .                                                          //
//                                                                             //
//    this = "this" .                                                          //
//                                                                             //
//    parent = "parent" .                                                      //
//                                                                             //
//    table = "table" .                                                        //
//                                                                             //
//    ident = (a-zA-Z_) .                                                      //
//                                                                             //
/////////////////////////////////////////////////////////////////////////////////
