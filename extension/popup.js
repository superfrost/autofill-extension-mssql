document.getElementById('fetchData').addEventListener('click', async () => {
  const response = await fetch('http://localhost:3000/data');
  const data = await response.json();
  console.log(data)
  document.getElementById('result').textContent = JSON.stringify(data, null, 2);
});
