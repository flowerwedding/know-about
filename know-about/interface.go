package main

/*

                                     参数名                            类型                          注释
接口1：/hello/first     前>>后       user_id                           string                       手机号
                                    password                          string                   用户输的验证码
                                  verification                        string                  前端生成的验证码
接口2：/hello/next      前>>后        name                             string                       昵称
                                  password                            string               表内其他默认为零，即空
接口3：/hello/last     前>>后      user_id                             string                       手机号
                                  password                            string

                                           参数                              类型                               注释
接口1：/cure/recommend                      无                                                           问题，该问题的赞最多的回答
接口2：/cure/follow                         无                                             关注的人的问题、回答、文章,如果没关注的人就推荐人和人
接口3：/cure/hostlist                       无                                                       返回前三名的问题，热度=关注+浏览
接口4：/cure/insert           前>>后      message                            string
                                          label                             string
                                          detail                            string                         问题详情、文章内容
                                          browse                            string                  问题browes == 0、文章browse == -1
                                           ni                               string                            ni == -1匿名
                                         picture                            string                    图片必须有值，没有就传默认图片
                              后>>前    消息里面要message和user_id
接口5：/cure/delete                       message                           string
接口6：/cure/fuzzysearch                 message                            string
                                           t                               string                    t -- 1是最后一次，t == 0是中间
                            后>>前          是最后一次返回所以相似问题，并且dynamic表内记录，相似问题有问题message，最佳回答
                                            不是最后一次返回相似问题和历史记录，相似问题有问题message

                                            参数名                      类型                         注释
接口1：/sun/findcomment     前>>后             w                       string                     什么问题
                                             pid                      string                   什么回答、问题
                                              f                       string            时间排序f == 1，默认排序f == 0
                            后>>前         评论所有message、子评论所有ChildrenMessage、评论者user_id、赞数、踩数、时间
                                           评论大于三条时默认排序还会返回精选评论三条
接口2：/sun/findanswer      前>>后            w                        string
                                             f                        string           时间排序f == 1，默认排序f == 0
                            后>>前         回答、回答者、赞数、踩数、评论数、时间
接口3：/sun/findall         前>>后            w                        string
                           后>>前         问题、问题详情、标签、关注数、浏览数、回答数、评论数
                                          精选回答三条f == 2
                                          第一个回答的关于作者：头像、昵称、回答数、文章数、关注者数
接口4：/sun/insert         前>>后           message                     string                     可以放表情包
	                                         pid                       string         pid == w写回答,pid == 0写问题的评论
	                                          w                        string
接口5：/sun/delete         前>>后          message                      string
接口6：/sun/update         前>>后          message                      string
                                              f                        string 关注人f == -1、关注问题f == 0、赞f == 1、踩f == 2、浏览问题f == 3、浏览人f == 4
                                             id                        string                        被赞的那个人
关注问题：问题表中的关注次数加一，动态表中的关注者有记录
浏览问题：动态表中没记录，问题表中有记录
关注人:动态表中关注者有记录
浏览人：动态表中没记录，测试表中有记录

                                       参数                              类型                           注释
接口1：/smile/update   前>>后           name                           string                            昵称
	                                  gender                          string                            性别
                                	background                        string                            背景
	                               headportrait                       string                            头像
	                                 introduce                        string                         一句话介绍
	                                  address                         string                           居住地
	                                  industry                        string                          所在行业
	                                 occupation                       string                          职业经历
	                                 education                        string                          教育经历
	                                  synopsis                        string                            简介
                                         t                            string                    查看t == 1,更新t == 0
更新必须都有值，尤其是图片，前端判断图片有误，没有就传默认图片
接口2：/smile/dynamic   前>>后         value                           string   动态、回答、问题、文章、专栏、想法、我关注的人、关注我的人、我关注的专栏、话题、问题、收藏
                       后>>前      文章：昵称，头像，文章，文章详情，赞
                                   我关注的人：对方的昵称，头像，回答数，文章数，关注着
                                   关注我的人：对方的昵称，头像，回答数，文章数，关注着
                                   动态：赞、踩、关注的问题、提问、回答
接口3：/smile/others   前>>后            无
                       后>>前       name,headportrait,background,count_huida,count_wenti,count_wenzhang,count_zhuanlan,count_xiangfa,count_woguanzhuderen,
                                    count_woguanzhudehuati,count_woguanzhudezhuanlan,count_woguanzhudewenti,count_woguanzhudeshoucangjia
接口4：/smile/logout

 */

/*

手机号及用户名
验证码前端生成

1.免密码登录
>>函数每次用接口时都调用>>判断是否为空
                  >>手机号如果为空，返回“请输入手机号”
                  >>验证码如果为空，返回“请输入验证码”
>>输入验证码，已有手机号(接口1)>>判断之前是否注册过
                             >>如果有，就登录，直接进入知乎
                             >>如果没有，就注册，并弹出注册页面
                                                >>输入注册用的用户名，密码(接口2)>>注册成功，进入知乎

2.密码登录
>>输入手机号或邮箱，密码(接口3)
           >>判断是否为空>>如果为空，返回“请输入手机号或邮箱”
                        >>如果不为空>>判断密码是否正确>>如果正确，进入知乎
                                                     >>如果不正确，返回“账号或密码错误”

3.第三方登录
1)微信
>>跳转小窗口
2)QQ
>>跳转小窗口
3)微博
>>进入微博登录页>>输入账号，密码(接口3)
               >>判断密码是否正确>>如果正确，进入知乎
                                >>如果不正确，返回“账号或密码错误”

 */