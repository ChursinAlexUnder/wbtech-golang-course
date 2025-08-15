// Когда вся страница html уже загружена
document.addEventListener("DOMContentLoaded", () => {

    const path = window.location.pathname;
    const searchButton = document.getElementById("searchButton");
    const searchInput = document.getElementById("searchInput");
    const errorMessage = document.getElementById("errorMessage")

    // Рекурсивная автоматическая функция обновления всех полей
    function updateElementsFromData(data) {
        function updateRecursive(obj) {
            Object.keys(obj).forEach(key => {
                // Если объект, то спускаемся внутрь, иначе если id, обновляем поле в html
                if (typeof obj[key] === "object" && obj[key] !== null) {
                    updateRecursive(obj[key]);
                } else {
                    const el = document.getElementById(key);
                    if (el) {
                        if (el.tagName === "TABLE" || el.closest("table")) {
                            return; 
                        }
                        el.textContent = obj[key];
                    }
                }
            });
        }
        updateRecursive(data);
    }

    // Запись полученных с сервера данных в таблицы
    function fillTable(tableId, array) {
        let table = document.getElementById(tableId);

        // Очищаем таблицу от старых данных
        while (table.rows.length > 1) {
            table.deleteRow(1);
        }

        // Берём первую строку (заголовки) и достаём id каждого th
        let headers = table.rows[0].querySelectorAll("th");
        let keys = Array.from(headers).map(th => th.id);

        // Создаем и заполняем
        array.forEach(item => {
            let row = table.insertRow();
            keys.forEach(key => {
            let cell = row.insertCell();
            cell.textContent = item[key];
            });
        });
    }
    
    // Функция для обращения на сервер по пути path
    function getDataFromServer(uuid) {
        fetch('/api/' + uuid)
            .then(response => response.json()) // Парсим ответ как JSON
            .then(data => {
                updateElementsFromData(data);
                fillTable("tableItems", data.items);
                fillTable("tableProduct", data.items);
            })
            .catch(error => console.error("Ошибка:", error));
    }

    // Проверка на правильность uid в URL
    function isCorrectURL(url) {
        return url.match(/^\/order\/([0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})$/i);
    }

    const match = isCorrectURL(path);
    if (match) {
        const uuid = match[1];
        getDataFromServer(uuid);
    }

    searchButton.addEventListener("click", () => {
        // валидация
        const newUrl = "/order/" + searchInput.value;
        if (isCorrectURL(newUrl)) {
            location.assign(newUrl);
        } else {
            errorMessage.style.display = 'block';
            searchInput.style.border = "0.2vw solid rgb(219, 23, 23)";
            searchInput.style.borderRight = "none";
            searchButton.style.border = "0.2vw solid rgb(219, 23, 23)";
            searchButton.style.borderLeft = "none";
        }
    });
});