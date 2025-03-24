<template>
  <!-- router-view отобразит текущее представление в зависимости от маршрута -->
  <router-view />
</template>

<script setup>

</script>

<style>
/* Заставляем весь документ иметь белый фон */
html, body, #app {
  background-color: white;
}

/* Можно задать шрифт по умолчанию и базовые стили */
body {
  font-family: Arial, sans-serif;
  color: #333;
}
</style>