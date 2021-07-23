The steps to deploy:

1. 安装python相关依赖
   
    ```
    sudo apt install python-cairo python-pip
    pip install cairosvg
    ```
    
2. 安装nodejs及npm

    ```
    sudo apt install nodejs npm
    ```

3. 使用npm更新nodejs(可选)

    ```
    sudo npm cache clean -f
    sudo npm install -g n
    sudo n stable
    ```

4. 进入mathjax目录

    ```
    npm install
    ```

5. 进入main.go所在目录

    ```
    go run main.go render.go s3communicate.go token.go
    ```

    

