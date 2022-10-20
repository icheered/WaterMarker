<script>
  import { SelectFolder, SelectFile, GetNumberOfFiles } from "../../wailsjs/go/main/App.js";
  import Button from "../lib/Button.svelte";

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

        fetch(window.location+ watermarkPath)
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

<div class="wrapper">
  <div>
    <Button text="Select Source Folder" callback={selectSourceFolderPath} />
    <div>
      {sourceFolderPath.match(/[^\/]+\/[^\/]+$/) ?? ""}
    </div>
    <div>
      {numberOfSourceFiles ? numberOfSourceFiles + " files found" : ""}
    </div>
  </div>
  <div>
    <Button text="Select Target Folder" callback={selectTargetFolderPath} />
    <div>
      {targetFolderPath.match(/[^\/]+\/[^\/]+$/) ?? ""}
    </div>
    <div>
      {numberOfTargetFiles ? numberOfTargetFiles + " files found" : ""}
    </div>
  </div>

  <div>
    <Button text="Select Watermark" callback={selectFilePath} />
    <div>
      {watermarkPath.match(/[^\/]+\/[^\/]+$/) ?? ""}
    </div>

    <img class="watermarkpreview" bind:this={watermarkpreviewImage} src="" alt="Watermark preview" />
  </div>
</div>

<style>
  img {
    max-width: 100px;
  }

  .watermarkpreview {
    min-height: 50px;
    min-width: 160px;
  }
</style>
