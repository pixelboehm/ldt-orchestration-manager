function toggleDescription(rowId) {
    var descriptionRow = document.getElementById('description-row-' + rowId);
    var button = document.getElementById('toggle-button-' + rowId);
    if (descriptionRow.style.display === 'none') {
      descriptionRow.style.display = 'table-row';
      button.innerHTML = 'Hide Description';
    } else {
      descriptionRow.style.display = 'none';
      button.innerHTML = 'Show Description';
    }
}