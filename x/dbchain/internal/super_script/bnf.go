package super_script

/////////////////////////////////////////////////////////////////////////////////
//                                                                             //
//                                     BNF                                     //
//                                                                             //
/////////////////////////////////////////////////////////////////////////////////
//                                                                             //
//    comparison = single_value "==" single_value .                            //
//               | single_value "in" multi_value .                             //
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
//    multi_value = table "." table_name { "." where } "." field               //
//                | "(" string_literal { "," string_literal } ")" .            //
//                                                                             //
//    table_name = ident                                                       //
//                                                                             //
//    where = "where" "(" field "==" single_value ")" .                        //
//                                                                             //
//    field = ident                                                            //
//                                                                             //
//    this = "this"                                                            //
//                                                                             //
//    parent = "parent"                                                        //
//                                                                             //
//    table = "table"                                                          //
//                                                                             //
//    ident = (a-zA-Z_)                                                        //
//                                                                             //
/////////////////////////////////////////////////////////////////////////////////
