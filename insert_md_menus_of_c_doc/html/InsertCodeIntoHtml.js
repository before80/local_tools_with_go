// 创建按钮元素
const button = document.createElement('button');
button.id = 'openUrlsButton';
button.textContent = '打开所有 URL';

// 创建脚本元素
const script = document.createElement('script');
script.textContent = `
    document.getElementById('openUrlsButton').addEventListener('click', function () {
        const mwGeshiDivs = document.querySelectorAll('div.mw-geshi');
        const allUrls = [];

        mwGeshiDivs.forEach((div) => {
            const aElements = div.querySelectorAll('a');
            aElements.forEach((a) => {
                const url = a.href;
                if (url) {
                    allUrls.push(url);
                }
            });
        });

        const uniqueUrls = [...new Set(allUrls)];
        console.log(uniqueUrls);
        let i = 0;
        uniqueUrls.forEach((url) => {
            // if (i < 20) { i++; } else { window.open(url, '_blank'); } 
            window.open(url, '_blank');           
        });
    });
`;

// 获取 body 元素
const body = document.body;

// 将按钮插入到 body 的最前面
if (body.firstChild) {
    body.insertBefore(button, body.firstChild);
} else {
    body.appendChild(button);
}

// 将脚本添加到页面的 body 元素中
body.appendChild(script);
