# Steps:

```sh
cd bot

sudo apt install python-cairo python-pip
pip install cairosvg

cd t2f-runtime/mathjax
npm install mathjax-full

cd ../..
chmod +x ./build-install.sh

tmux
./build-install.sh
```

然后按下`ctrl b`, 再按下`d`
