const barcodeInput = document.querySelector("#barcode-values");
const KEYCODE_ENTER = 13;
const KEYCODE_CTRL_ENTER = 10;

window.addEventListener("load", function () {
  const urlSearchParams = new URLSearchParams(window.location.search);
  barcodeInput.value = urlSearchParams.getAll("b[]").join("\n");

  const selectedType = urlSearchParams.get('type');
  if (!selectedType) return;

  const barcodeType = document.querySelector(`input[name="type"][value="${selectedType}"]`);
  barcodeType.checked = true;
});

function generateBarcodes() {
  const barcodeInputValue = barcodeInput.value.trim();
  if (barcodeInputValue.length == 0) return;

  const barcodeType = document.querySelector('input[name="type"]:checked');
  const baseURL = window.location.protocol + "//" + window.location.host + "/";

  const barcodeValues = barcodeInputValue.split("\n");
  const params = ["type=" + barcodeType.value];

  for (var i = 0; i < barcodeValues.length; i++) {
    params.push("b[]=" + encodeURIComponent(barcodeValues[i]));
  }

  window.location.href = baseURL + "?" + params.join("&");
}

document.querySelector("#generate-barcodes").addEventListener("click", generateBarcodes);

barcodeInput.addEventListener("keydown", function (e) {
  if ((e.ctrlKey || e.metaKey) && (e.keyCode === KEYCODE_ENTER || e.keyCode === KEYCODE_CTRL_ENTER)) {
    generateBarcodes();
  }
});

const pdfUnitSelector = document.querySelector('#pdf-unit');
const dimensionInputs = document.querySelectorAll('.dimension-input');
pdfUnitSelector.addEventListener('change', function(e) {
  const showDimensionInputs = e.target.value !== 'auto';
  dimensionInputs.forEach(input => showDimensionInputs ? input.style.display = 'inherit' : input.style.display = 'none');
})
