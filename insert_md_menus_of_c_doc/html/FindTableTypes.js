/*
* 通过表格找出所有类型、宏、函数
* */

// 获取所有类名为 t-dsc-begin 的表格元素
const tables = document.querySelectorAll('table.t-dsc-begin');
const firstColumnContents = [];

// 遍历每个表格
tables.forEach(table => {
    // 获取表格的所有行
    const rows = table.querySelectorAll('tr.t-dsc');
    rows.forEach(row => {
        // 获取每行的第一列
        const firstCell = row.querySelector('td .t-dsc-member-div div:first-child');
        if (firstCell) {
            // 获取第一列中的 span 元素文本内容
            const spanElements = firstCell.querySelectorAll('span.t-lines span');
            spanElements.forEach(spanElement => {
                // 将内容添加到数组中
                firstColumnContents.push(spanElement.textContent.trim());
            })
        }
    });
});

// 将数组中的内容用空格连接成字符串
const result = firstColumnContents.join(' ');
console.log(result);
    