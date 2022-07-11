/*********************************************
    AUTO GENERATED BY TOOLS, DO NOT EDIT
*********************************************/
using System;
using System.Collections.Generic;

namespace Senle.Config
{
    public class t_error
    {
        public class t_error_item
        {
            /// <summary>
            /// 键
            /// </summary>
            public string key { get; internal set; }
                        
            /// <summary>
            /// 值
            /// </summary>
            public int code { get; internal set; }
                        
            /// <summary>
            /// 描述
            /// </summary>
            public string desc { get; internal set; }
                        


        }



        private static Dictionary< content = new Dictionary<();
        
        public static Dictionary< get_all()
        {
            return content;
        }

        


        static t_error()
        {
            content = __data__();
        }
        private static Dictionary< __data__()
        {
            return new Dictionary<()
            {
                {"SEVER_BUSY", new ()
                    {
                        key = "SEVER_BUSY",
                        code = 101,
                        desc = "服务器发生错误或者繁忙",
                    }
                },
                {101, new ()
                    {
                        key = "SEVER_BUSY",
                        code = 101,
                        desc = "服务器发生错误或者繁忙",
                    }
                },
                {"WC_SERVER_BUSY", new ()
                    {
                        key = "WC_SERVER_BUSY",
                        code = 102,
                        desc = "微信服务器异常",
                    }
                },
                {102, new ()
                    {
                        key = "WC_SERVER_BUSY",
                        code = 102,
                        desc = "微信服务器异常",
                    }
                },
                {"WC_INVALID_CODE", new ()
                    {
                        key = "WC_INVALID_CODE",
                        code = 103,
                        desc = "非法登录码",
                    }
                },
                {103, new ()
                    {
                        key = "WC_INVALID_CODE",
                        code = 103,
                        desc = "非法登录码",
                    }
                },
                {"WC_TOO_MANY_LOGIN_TIMES", new ()
                    {
                        key = "WC_TOO_MANY_LOGIN_TIMES",
                        code = 104,
                        desc = "登录太频繁",
                    }
                },
                {104, new ()
                    {
                        key = "WC_TOO_MANY_LOGIN_TIMES",
                        code = 104,
                        desc = "登录太频繁",
                    }
                },
                {"INVALID_PARAMS", new ()
                    {
                        key = "INVALID_PARAMS",
                        code = 105,
                        desc = "参数不正确",
                    }
                },
                {105, new ()
                    {
                        key = "INVALID_PARAMS",
                        code = 105,
                        desc = "参数不正确",
                    }
                },
                {"INVALID_TOKEN", new ()
                    {
                        key = "INVALID_TOKEN",
                        code = 106,
                        desc = "密钥不正确",
                    }
                },
                {106, new ()
                    {
                        key = "INVALID_TOKEN",
                        code = 106,
                        desc = "密钥不正确",
                    }
                },
            }

;
        }

    }
}

