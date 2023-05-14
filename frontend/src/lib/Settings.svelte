<script>
  export let watermarkOpacity, watermarkScale, watermarkPosition, changedSettings, parallelProcesses;

  function valuechanged() {
    changedSettings = true;
  }
</script>

<div>
  <div class="row">
    <div class="text">Opacity</div>
    <input class="textfield" type="text" bind:value={watermarkOpacity} on:change={valuechanged} min="0" max="100" />
  </div>
  <div>
    <input class="slider" type="range" bind:value={watermarkOpacity} on:change={valuechanged} min="0" max="100" />
  </div>

  <div class="row">
    <div class="text">Scale</div>
    <input class="textfield" type="text" bind:value={watermarkScale} on:change={valuechanged} min="0" max="100" />
  </div>
  <div>
    <input class="slider" type="range" bind:value={watermarkScale} on:change={valuechanged} min="0" max="100" />
  </div>

  <div class="row">
    <div class="text">Parallel processes</div>
    <input class="textfield" type="text" bind:value={parallelProcesses} min="1" max="64" />
  </div>
  <div>
    <input class="slider" type="range" bind:value={parallelProcesses} min="1" max="64" />
  </div>


  <div>Watermark location</div>
  <select bind:value={watermarkPosition} on:change={valuechanged}>
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
  .textfield {
    width: 25px;
    margin-left: 10px;
  }
  .text {
    padding-top: 3px;
  }

  @media screen and (-webkit-min-device-pixel-ratio: 0) {
    input[type="range"] {
      overflow: hidden;
      -webkit-appearance: none;
      background-color: white;
      border-radius: 5px;
    }

    input[type="range"]::-webkit-slider-runnable-track {
      height: 10px;
      -webkit-appearance: none;
      color: #13bba4;
      margin-top: -1px;
    }

    input[type="range"]::-webkit-slider-thumb {
      width: 6px;
      -webkit-appearance: none;
      height: 10px;
      cursor: ew-resize;
      background: black;
      box-shadow: -80px 0 0 80px #1caeee;
    }
  }
</style>
