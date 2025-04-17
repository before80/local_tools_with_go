/*
* 找出“表格”中的所有链接并按照指定的起始和所需个数进行打开
* */
function CreateOpenUrlButtonFromTable(startNum, needNum) {
    window.myStartNum = startNum;
    window.myEndNum = startNum + needNum;
    if (document.getElementById('openUrlsButton')) {
        return
    }

    // 创建按钮元素
    const button = document.createElement('button');
    button.id = 'openUrlsButton';
    button.textContent = '打开所有 URL';

    // 创建脚本元素
    const script = document.createElement('script');
    script.textContent = `
    document.getElementById('openUrlsButton').addEventListener('click', function () {
        const tables = document.querySelectorAll('table.t-dsc-begin');
        const allUrls = [];

        tables.forEach((table) => {
            const rows = table.querySelectorAll('tr.t-dsc');
            rows.forEach(row => {
                // 获取每行的第一列
                const firstCell = row.querySelector('td .t-dsc-member-div div:first-child');
                if (firstCell) {
                    const aElement = firstCell.querySelector('a');
                    const url = aElement.href;
                    if (url) {
                        allUrls.push(url);
                    }
                }
            })
        });

        const uniqueUrls = [...new Set(allUrls)];
        console.log(uniqueUrls);
        let i = 0;
        uniqueUrls.forEach((url) => {
            if (i >= myStartNum && i < myEndNum) { window.open(url, '_blank');}
            i++;
            // window.open(url, '_blank');
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
}

CreateOpenUrlButtonFromTable(0, 10)

