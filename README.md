selpg
====
使用golang开发[Linux命令行实用程序](https://www.ibm.com/developerworks/cn/linux/shell/clutil/index.html)中的selpg
-------
* 代码结构
  
    根据文档提供的c语言变换成golang实现，主体大致有三个函数实现：
  
  ```
  func main()
  func process_args(ac int, args []string, psa *selpg_args)
  func process_input(sa selpg_args)
  ```
  main()函数声明所需变量，并调用函数process_args()解析命令行参数，process_input()函数根据解析后的参数选择所需页写至制定目的地。
  
  ```
  type selpg_args struct {
	start_page  int
	end_page    int
	in_filename string
	page_len    int  
  page_type rune
	print_dest  string
  }
  ```
  这里定义了selpg_args结构，存储传递进来的命令行参数。
* 测试
  
  
  已有文件inputfile，outputfile。
  其中inputfile内容如下,outputfile文件为空
  ```
  test1
  test2
  test3
  test4
  ```
  
  执行下列命令：
  
  ```
  $ go run selpg.go -s1 -e1 inputfile
  ```
  
  
  该命令将inputfile的第一页写至标准输出即屏幕，结果验证正确。
  ```
  $ go run selpg.go -s1 -e1 < inputfile
  ```
  
  该命令使selpg读取标准输入，而标准输入已被内核重定向为input file，所以此处仍然在屏幕上显示inputfile第一页的内容，结果验证正确。
  

  ```
  $ go run selpg.go -s1 -e1 inputfile| selpg.go -s10 -e20
  ```
  
  该命令读取inputfile第一页作为输入，将第10页到第20页写到标准屏幕，这里显然开始页大于总页数，会返回错误显示开始页大于总页数。
  ```
  $ go run selpg.go -s1 -e1 inputfile >outputfile
  ```
  该命令将input file的第一页作为标准输入，标准输出被内核重定向为outputfile，验证结果outputfile与inputfile第一页内容相同，结果正确。
  ```
  $ go run selpg.go -s1 -s2 inputfile 2>errorfile
  ```
  该命令读取inputfile的前两页至标准输出，错误消息被重定向至errorfile，此时没有errorfile，会自动生成该文件并将错误信息写入，结果验证正确。
  ```
  $ go run selpg.go -s1 -s2 inputfile >outputfile 2>errorfile 
  ```
  该命令读取inputfile的前两页作为标准输出，标准输出被重定向为outputfile，错误消息被重定向为errorfile，执行后会观察到errorfile显示结束页大于总页数的错误，结果验证正确。
  ```
  $ go run selpg.go -s2 -e2 -l2 inputfile
  ```
  该命令将页长设置为2，把输入根据页长输出第二页，结果验证正确。
  ```
  $ go run selpg.go -s2 -e2 -f inputfile
  ```
  该命令使页由换页符界定，输出第二页。
  
* 实验总结
  
  
  这次实验主要是学习golang语言，有很大收获，其中关于设备输出的部分没有写到，有不足之处，日后还要多多学习。
  
