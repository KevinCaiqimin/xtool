/*********************************************
    AUTO GENERATED BY TOOLS, DO NOT EDIT
*********************************************/
using System;
using System.Collections.Generic;

namespace Senle.Config
{
    public class lv_cfg
    {
        public class lv_cfg_item
        {
            /// <summary>
            /// 宠物ID
            /// </summary>
            public string id { get; internal set; }
                        
            /// <summary>
            /// 等级
            /// </summary>
            public int lv { get; internal set; }
                        
            /// <summary>
            /// 需要经验
            /// </summary>
            public int needExp { get; internal set; }
                        
            /// <summary>
            /// 猫猫头像
            /// </summary>
            public string headImg { get; internal set; }
                        
            /// <summary>
            /// 展示图片
            /// </summary>
            public string img { get; internal set; }
                        
            /// <summary>
            /// 动画
            /// </summary>
            public string animate { get; internal set; }
                        


        }



        private static Dictionary<string,Dictionary<int,lv_cfg_item>> content = new Dictionary<string,Dictionary<int,lv_cfg_item>>();
        
        public static Dictionary<string,Dictionary<int,lv_cfg_item>> get_all()
        {
            return content;
        }

        public static Dictionary<int,lv_cfg_item> get_conf(string id)
        {
            var dic1 = content;
            if (!dic1.ContainsKey(id))
            {
                return null;
            }
            


            return dic1[id];
        }
        public static lv_cfg_item get_conf(string id, int lv)
        {
            var dic1 = content;
            if (!dic1.ContainsKey(id))
            {
                return null;
            }
            
            var dic2 = dic1[id];
            if (!dic2.ContainsKey(lv))
            {
                return null;
            }
            


            return dic2[lv];
        }



        static lv_cfg()
        {
            content = __data__();
        }
        private static Dictionary<string,Dictionary<int,lv_cfg_item>> __data__()
        {
            return new Dictionary<string,Dictionary<int,lv_cfg_item>>()
            {
                {"cat_huahua", new Dictionary<int,lv_cfg_item>()
                    {
                        {1, new lv_cfg_item()
                            {
                                id = "cat_huahua",
                                lv = 1,
                                needExp = 10,
                                headImg = "defaultImg.png",
                                img = "defaultImg.png",
                                animate = "",
                            }
                        },
                        {2, new lv_cfg_item()
                            {
                                id = "cat_huahua",
                                lv = 2,
                                needExp = 20,
                                headImg = "defaultImg.png",
                                img = "defaultImg.png",
                                animate = "",
                            }
                        },
                        {3, new lv_cfg_item()
                            {
                                id = "cat_huahua",
                                lv = 3,
                                needExp = 30,
                                headImg = "defaultImg.png",
                                img = "defaultImg.png",
                                animate = "",
                            }
                        },
                        {4, new lv_cfg_item()
                            {
                                id = "cat_huahua",
                                lv = 4,
                                needExp = 40,
                                headImg = "defaultImg.png",
                                img = "defaultImg.png",
                                animate = "",
                            }
                        },
                        {5, new lv_cfg_item()
                            {
                                id = "cat_huahua",
                                lv = 5,
                                needExp = 50,
                                headImg = "defaultImg.png",
                                img = "defaultImg.png",
                                animate = "",
                            }
                        },
                        {6, new lv_cfg_item()
                            {
                                id = "cat_huahua",
                                lv = 6,
                                needExp = 60,
                                headImg = "defaultImg.png",
                                img = "defaultImg.png",
                                animate = "",
                            }
                        },
                        {7, new lv_cfg_item()
                            {
                                id = "cat_huahua",
                                lv = 7,
                                needExp = 70,
                                headImg = "defaultImg.png",
                                img = "defaultImg.png",
                                animate = "",
                            }
                        },
                        {8, new lv_cfg_item()
                            {
                                id = "cat_huahua",
                                lv = 8,
                                needExp = 80,
                                headImg = "defaultImg.png",
                                img = "defaultImg.png",
                                animate = "",
                            }
                        },
                        {9, new lv_cfg_item()
                            {
                                id = "cat_huahua",
                                lv = 9,
                                needExp = 90,
                                headImg = "defaultImg.png",
                                img = "defaultImg.png",
                                animate = "",
                            }
                        },
                        {10, new lv_cfg_item()
                            {
                                id = "cat_huahua",
                                lv = 10,
                                needExp = 100,
                                headImg = "defaultImg.png",
                                img = "defaultImg.png",
                                animate = "",
                            }
                        },
                    }
                },
                {"cat_miaomiao", new Dictionary<int,lv_cfg_item>()
                    {
                        {1, new lv_cfg_item()
                            {
                                id = "cat_miaomiao",
                                lv = 1,
                                needExp = 10,
                                headImg = "defaultImg.png",
                                img = "defaultImg.png",
                                animate = "",
                            }
                        },
                        {2, new lv_cfg_item()
                            {
                                id = "cat_miaomiao",
                                lv = 2,
                                needExp = 20,
                                headImg = "defaultImg.png",
                                img = "defaultImg.png",
                                animate = "",
                            }
                        },
                        {3, new lv_cfg_item()
                            {
                                id = "cat_miaomiao",
                                lv = 3,
                                needExp = 30,
                                headImg = "defaultImg.png",
                                img = "defaultImg.png",
                                animate = "",
                            }
                        },
                        {4, new lv_cfg_item()
                            {
                                id = "cat_miaomiao",
                                lv = 4,
                                needExp = 40,
                                headImg = "defaultImg.png",
                                img = "defaultImg.png",
                                animate = "",
                            }
                        },
                        {5, new lv_cfg_item()
                            {
                                id = "cat_miaomiao",
                                lv = 5,
                                needExp = 50,
                                headImg = "defaultImg.png",
                                img = "defaultImg.png",
                                animate = "",
                            }
                        },
                        {6, new lv_cfg_item()
                            {
                                id = "cat_miaomiao",
                                lv = 6,
                                needExp = 60,
                                headImg = "defaultImg.png",
                                img = "defaultImg.png",
                                animate = "",
                            }
                        },
                        {7, new lv_cfg_item()
                            {
                                id = "cat_miaomiao",
                                lv = 7,
                                needExp = 70,
                                headImg = "defaultImg.png",
                                img = "defaultImg.png",
                                animate = "",
                            }
                        },
                        {8, new lv_cfg_item()
                            {
                                id = "cat_miaomiao",
                                lv = 8,
                                needExp = 80,
                                headImg = "defaultImg.png",
                                img = "defaultImg.png",
                                animate = "",
                            }
                        },
                        {9, new lv_cfg_item()
                            {
                                id = "cat_miaomiao",
                                lv = 9,
                                needExp = 90,
                                headImg = "defaultImg.png",
                                img = "defaultImg.png",
                                animate = "",
                            }
                        },
                        {10, new lv_cfg_item()
                            {
                                id = "cat_miaomiao",
                                lv = 10,
                                needExp = 100,
                                headImg = "defaultImg.png",
                                img = "defaultImg.png",
                                animate = "",
                            }
                        },
                    }
                },
                {"cat_tangyuan", new Dictionary<int,lv_cfg_item>()
                    {
                        {1, new lv_cfg_item()
                            {
                                id = "cat_tangyuan",
                                lv = 1,
                                needExp = 10,
                                headImg = "defaultImg.png",
                                img = "defaultImg.png",
                                animate = "",
                            }
                        },
                        {2, new lv_cfg_item()
                            {
                                id = "cat_tangyuan",
                                lv = 2,
                                needExp = 20,
                                headImg = "defaultImg.png",
                                img = "defaultImg.png",
                                animate = "",
                            }
                        },
                        {3, new lv_cfg_item()
                            {
                                id = "cat_tangyuan",
                                lv = 3,
                                needExp = 30,
                                headImg = "defaultImg.png",
                                img = "defaultImg.png",
                                animate = "",
                            }
                        },
                        {4, new lv_cfg_item()
                            {
                                id = "cat_tangyuan",
                                lv = 4,
                                needExp = 40,
                                headImg = "defaultImg.png",
                                img = "defaultImg.png",
                                animate = "",
                            }
                        },
                        {5, new lv_cfg_item()
                            {
                                id = "cat_tangyuan",
                                lv = 5,
                                needExp = 50,
                                headImg = "defaultImg.png",
                                img = "defaultImg.png",
                                animate = "",
                            }
                        },
                        {6, new lv_cfg_item()
                            {
                                id = "cat_tangyuan",
                                lv = 6,
                                needExp = 60,
                                headImg = "defaultImg.png",
                                img = "defaultImg.png",
                                animate = "",
                            }
                        },
                        {7, new lv_cfg_item()
                            {
                                id = "cat_tangyuan",
                                lv = 7,
                                needExp = 70,
                                headImg = "defaultImg.png",
                                img = "defaultImg.png",
                                animate = "",
                            }
                        },
                        {8, new lv_cfg_item()
                            {
                                id = "cat_tangyuan",
                                lv = 8,
                                needExp = 80,
                                headImg = "defaultImg.png",
                                img = "defaultImg.png",
                                animate = "",
                            }
                        },
                        {9, new lv_cfg_item()
                            {
                                id = "cat_tangyuan",
                                lv = 9,
                                needExp = 90,
                                headImg = "defaultImg.png",
                                img = "defaultImg.png",
                                animate = "",
                            }
                        },
                        {10, new lv_cfg_item()
                            {
                                id = "cat_tangyuan",
                                lv = 10,
                                needExp = 100,
                                headImg = "defaultImg.png",
                                img = "defaultImg.png",
                                animate = "",
                            }
                        },
                    }
                },
            }

;
        }

    }
}
