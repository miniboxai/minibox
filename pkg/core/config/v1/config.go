package v1

import (
	"net/url"

	"minibox.ai/pkg/core/node"
	"minibox.ai/pkg/core/option"
)

// ---
// train:
//   framework: torch@2.1.7
//
// ---
// train:
//   framework:
//     name: torch
//     version: 2.1.7
//
// ---
// train:
//   framework: tensorflow@1.8.1
//
// ---
// train:
//   framework:
//     name: torch
//     version: 2.1.7
//
//
// train:
//   framework:
//     name: tensorflow
//     version: 1.8.0
//     sense: py2

// ---
// train:
//   build: ./  # 默认是字符串，指定编译目录 , Dockerfile 是默认文件名
//
// ---
// train:
//   build:
//     context: ./
//     dockerfile: build.docker # 指定编译文件名 build.docker
//
// ---
// train:
//   dirs:
//     data: ./another-data   # 指定 datadir 的目录名，默认是 ./data 这里，注：只能使用相对路径
//                            # 这个设置会影响，容器的 env["DATADIR"] , 即 DATADIR 环境变量
// 	                          # 可以在源代码码中，获得正确的，数据加入的路径
//     log: ./middle-logger   # 训练数据中间缓存目录，默认是 ./logdir
//     model:                 # 模型生成后的目录
//
// ---
// train:
//   datasets:                          # Dataset 数据列表
//     - mnist 	                        # 指向 minibox.ai/datasets/mnist
//                                      # 数据包下载到 datadir 的位置, 同时以数据集的名称命名的子目录
//                                      # 像：./data/mnist
//
//     - mnist@v1.2.3    				# 下载 v1.2.3 的版本，映射成 mnist-v1.2.3
//     - bob/mnist       				# 下载用户名字空间的数据集, 映射成 bob/mnist
//                                      # minibox.ai/bob/datasets/mnist
//
//     - mycompany.co/mnist.tar.gz  			# 下载 mycompany.co 域名下的 /mnist.tar.gz
//     - localhost:5555/mnist.tar.gz 			# 也可以指定端口, 并不需要指定 http[s]
//
//     # 数据集的名称不能相同， 否则映射就会失败，默认情况下，数据集都会转换成，相对应的名称
//     # 只有带域名地址的数据源，可是以文件结尾，其它的默认为，minibox 数据集合上的，数据集
//     # 数据集是若干个，数据文件的集合，原则上：数据文件，是一包含大量，样本与标签的集合，
//     # 我们不推荐，直接获得样本与标签，而且整合了大量的样本与标签的数据文件，可以是 tensorflow
//     # tfrecords, 或 hdf5，或者是 tar.gz 等压缩文件,
//     # 在容器中可以解压这些数据库，但是并不会持久存储，需要注意以免丢失数据
//     # 数据集映射到容器中，是只读的，你可以通过平台上传，或 mini data cp 到存储空间中去
//
// ---
// train:
//   model: trainer.task        # 模型载入器, trainer.task 对应的应该是，语言的包名字空间
//
//     entry: trainer/task.py   # 入口文件
//     args:					# 参数
//       - start
//     launch: /usr/bin/python  # 载入程序

// ---
// train:                       # 还有一种更直观的启动方式
//   cmd:
//     - python
//     - trainer/task.py
//     - --logdir  ./logdir
//     - --datadir ./data
//
// ---
// train:
//   env:
//     - http_proxy=127.0.0.1:1087  # 设定环境变量
//
// ---
// train:
//   ports:                         # 暴露端口，local 模式才有效
//     - 3000:4000
//
// ---
// train:
//   notebook: true                 # 开启 notebook
//
//   notebook:
//     image: torch-notebook/notebook:latest # 指定启动镜像
//
//     build:                     # build/image 二选一，定义了 build, image 的配置就失效了
//       context: ./
//       dockerfile: build.docker # 指定编译文件名 build.docker
//
//     cmd:
//       - jupyter
//       - notebook
//     context: /notebook					 # notebook 的工作目录
//     dir: ./notebooks                      # notebooks 文件来源
//     port: 8888							 # 暴露端口 (only local)
//
// ---
// train:
//   dashboard: true                # 开启 dashboard
//
type Config struct {
	opts      *option.ConfigOpt
	nodes     map[string]node.Noder
	Framework Framework // 训练的软件环境描述，反映出语言与框架，让 Backend 决定使用那种 Image 来执行
	Build     Build     // 自定义用户的 Image 环境，来至于 docker-compose 的 build 格式。
	Workdir   string    // 工作目录，映射到容器执行目录
	Dirs      Dirs      // 目录
	DataSets  []url.URL // 数据集会被映射到 datadir 下面
	Model     Model     // 模型启动器
	Cmd       []string  // 直接控制启动命令
	Env       []Env     // 环境变量设置
	Ports     []PortMap // 端口映射
	Notebook  Notebook  // Jupyter Notebook 服务器
	Dashboard Dashboard // Tensorboard 服务器

	hasFrame bool
	hasBuild bool
}

type Component struct {
	Image   string
	Build   Build
	Port    int
	Context string
	Cmd     []string
	Env     []Env
	links   []Link
	ports   []PortMap
}

type Require struct {
	Name    string
	Version string
}

type Arg struct {
	Name  string
	Value string
}

type MapArgs map[string]Arg

type Build struct {
	Context    string
	Dockerfile string
	Args       MapArgs
}

// 目录选项, 会在容器中，反应在 ENV["DATADIR"], ENV["LOGDIR"], ENV["MODELDIR"]
type Dirs struct {
	Data  string
	Logs  string
	Model string
}

type Env struct {
	Name  string
	Value string
}

type Link struct {
	ContainerName string
}

type PortMap struct {
	Src  int
	Dest int
}

type Model struct {
}

type Notebook struct {
	Component
	Dir         string
	PrimaryFile string
}

type Dashboard struct {
	Component
}
