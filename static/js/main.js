document.addEventListener('DOMContentLoaded', () => {
    // Получаем элементы
    const dropdownButton = document.getElementById('dropdownButton');
    const dropdownMenu = document.getElementById('dropdownMenu');

    // Флаг для отслеживания состояния
    let isMenuOpen = false;

    // Обработчик клика по кнопке
    dropdownButton.addEventListener('click', (event) => {
        event.stopPropagation(); // Останавливаем всплытие события
        isMenuOpen = !isMenuOpen;
        dropdownMenu.classList.toggle('hidden', !isMenuOpen);
    });

    // Закрытие меню при клике вне его области
    document.addEventListener('click', () => {
        isMenuOpen = false;
        dropdownMenu.classList.add('hidden');
    });
});