<script>
  import {
    SelectFolder,
    SelectFile,
    ProcessImages,
    GetNumberOfFiles,
    FetchPreview,
  } from "../../wailsjs/go/main/App.js";

  export let sourceFolderPath,
    targetFolderPath,
    watermarkPath,
    numberOfSourceFiles,
    numberOfTargetFiles,
    watermarkpreviewImage;

  function selectSourceFolderPath() {
    SelectFolder().then((result) => {
      if (result.length) {
        sourceFolderPath = result;
        GetNumberOfFiles(sourceFolderPath).then((result) => {
          numberOfSourceFiles = result;
        });
      }
    });
  }
  function selectTargetFolderPath() {
    SelectFolder().then((result) => {
      if (result.length) {
        targetFolderPath = result;
        GetNumberOfFiles(targetFolderPath).then((result) => {
          numberOfTargetFiles = result;
        });
      }
    });
  }

  function selectFilePath() {
    SelectFile().then((result) => {
      if (result.length) {
        watermarkPath = result;

        fetch(watermarkPath)
          .then((response) => response.blob())
          .then((blob) => {
            const reader = new FileReader();
            reader.addEventListener("load", function () {
              watermarkpreviewImage.setAttribute("src", reader.result);
            });
            reader.readAsDataURL(blob);
          });
      }
    });
  }
</script>

<div>
  <div>
    <button on:click={selectSourceFolderPath}>Select image folder</button>
    <div>
      {sourceFolderPath.match(/[^\/]+\/[^\/]+$/) ?? ""}
    </div>
    <div>
      {numberOfSourceFiles ? numberOfSourceFiles + " files found" : ""}
    </div>
  </div>
  <div>
    <button on:click={selectTargetFolderPath}>Select output folder</button>
    <div>
      {targetFolderPath.match(/[^\/]+\/[^\/]+$/) ?? ""}
    </div>
  </div>

  <div>
    <button on:click={selectFilePath}>Select Watermark</button>
    <div>
      {watermarkPath.match(/[^\/]+\/[^\/]+$/) ?? ""}
    </div>

    <img bind:this={watermarkpreviewImage} src="" alt="Watermark preview" />
  </div>
</div>

<style>
  img {
    max-width: 100px;
  }
</style>
