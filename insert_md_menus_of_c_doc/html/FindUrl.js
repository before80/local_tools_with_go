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

function openUrlsSequentially(urls, index = 0) {
    if (index < urls.length) {
        window.open(urls[index], '_blank');
        setTimeout(() => {
            openUrlsSequentially(urls, index + 1);
        }, 1000); // 每隔 1 秒打开一个新标签页
    }
}

openUrlsSequentially(uniqueUrls);