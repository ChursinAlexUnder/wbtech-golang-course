// Когда вся страница html уже загружена
document.addEventListener("DOMContentLoaded", () => {

    // Получение элементов страниц
    const searchInput = document.getElementById('searchInput');
// ...
    
    // Функция для обращения на сервер по пути path
    function getDataFromServer(uuid) {
        fetch('/api/' + uuid)
            .then(response => response.json()) // Парсим ответ как JSON
            .then(data => {
                // console.log("Ответ сервера:", data);
                // alert(`Статус: ${data.status}, Число: ${data.number}`);
            })
            .catch(error => console.error("Ошибка:", error));
    }

    // Проверка на ручной переход по URL
    const path = window.location.pathname;
    const match = path.match(/^\/order\/([0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})$/i);
    if (match) {
        const uuid = match[1];
        getDataFromServer(uuid);
    }



});