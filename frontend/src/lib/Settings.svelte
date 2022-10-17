<script>
  export let watermarkOpacity, watermarkScale, watermarkPosition;

  function changethis(e) {
    console.log("Changing this");
    console.log(e.target.value);
    watermarkPosition = e.target.value;
  }
</script>

<div>
  <div class="row">
    <div>Opacity</div>
    <input class="textfield" type="text" bind:value={watermarkOpacity} min="0" max="100" />
  </div>
  <div>
    <input type="range" bind:value={watermarkOpacity} min="0" max="100" />
  </div>

  <div class="row">
    <div>Scale</div>
    <input class="textfield" type="text" bind:value={watermarkScale} min="0" max="100" />
  </div>
  <div>
    <input type="range" bind:value={watermarkScale} min="0" max="100" />
  </div>

  <select bind:value={watermarkPosition} on:change={changethis}>
    <option value="bottom-right">Bottom-right</option>
    <option value="bottom-left">Bottom-left</option>
    <option value="top-right">Top-right</option>
    <option value="top-left">Top-left</option>
    <option value="center">Center</option>
  </select>
</div>

<style>
  .row {
    display: flex;
    flex-direction: row;
    justify-content: center;
  }
  .col {
    display: flex;
    flex-direction: column;
  }
  .textfield {
    width: 25px;
  }
</style>
