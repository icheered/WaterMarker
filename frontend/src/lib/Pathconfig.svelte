<script>
  import {
    SelectFolder,
    SelectFile,
    GetNumberOfFiles,
  } from "../../wailsjs/go/main/App.js";
  import Button from "../lib/Button.svelte";

  export let sourceFolderPath,
    targetFolderPath,
    watermarkPath,
    numberOfSourceFiles,
    numberOfTargetFiles,
    watermarkpreviewImage,
    changedSettings,
    showImagePreview;

  let showWatermarkPreview = false;

  function selectSourceFolderPath() {
    SelectFolder().then((result) => {
      if (result.length) {
        sourceFolderPath = result;
        changedSettings = true;
        showImagePreview = false;
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
          if (numberOfTargetFiles) {
            if (
              confirm(
                "The output folder is not empty, files may be overwritten. Continue?"
              )
            ) {
              // Save it!
              console.log("Target folder is not empty but continuing anyways");
            } else {
              // Do nothing!
              console.log("Not continuing");
              targetFolderPath = ""
              numberOfTargetFiles = 0
            }
          }
        });
      }
    });
  }

  function selectFilePath() {
    SelectFile().then((result) => {
      if (result.length) {
        watermarkPath = result;

        const fetchpath = navigator.platform.includes("Linux")
          ? watermarkPath
          : window.location + watermarkPath;
        fetch(fetchpath)
          .then((response) => response.blob())
          .then((blob) => {
            showWatermarkPreview = true;
            showImagePreview = false;
            const reader = new FileReader();
            reader.addEventListener("load", function () {
              changedSettings = true;
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
      {sourceFolderPath.match(/[^\/|\\]+[\/|\\][^(\/|\\)]+$/) ?? ""}
    </div>
    <div>
      {numberOfSourceFiles ? numberOfSourceFiles + " files found" : ""}
    </div>
  </div>
  <div>
    <Button text="Select Target Folder" callback={selectTargetFolderPath} />
    <div>
      {targetFolderPath.match(/[^\/|\\]+[\/|\\][^(\/|\\)]+$/) ?? ""}
    </div>
    <div>
      {numberOfTargetFiles ? numberOfTargetFiles + " files found" : ""}
    </div>
  </div>

  <div>
    <Button text="Select Watermark" callback={selectFilePath} />
    <div>
      {watermarkPath.match(/[^\/|\\]+[\/|\\][^(\/|\\)]+$/) ?? ""}
    </div>

    {#if showWatermarkPreview == true}
      <img
        class="watermarkpreview"
        bind:this={watermarkpreviewImage}
        src=""
        alt="Watermark preview"
      />
    {/if}
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
