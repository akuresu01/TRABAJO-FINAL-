document.getElementById('co2Form').addEventListener('submit', function(event) {
    event.preventDefault();
    const formData = new FormData(this);
    const data = Object.fromEntries(formData.entries());

    fetch('http://localhost:8080/calculate', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(data)
    })
    .then(response => response.json())
    .then(data => {
        document.getElementById('result').textContent = 'Total de CO2 gastado: ' + data.totalCO2.toFixed(2) + ' toneladas';
    })
    .catch((error) => {
        console.error('Error:', error);
        document.getElementById('result').textContent = 'Error en el cálculo. Por favor, inténtelo de nuevo.';
    });
});
