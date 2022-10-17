<script>
  import { ProcessImages, FetchPreview } from "../wailsjs/go/main/App.js";
  import Header from "./lib/Header.svelte";
  import Pathconfig from "./lib/Pathconfig.svelte";
  import Settings from "./lib/Settings.svelte";

  // CONSTANTS
  const title = "WaterMarker";
  const version = "1.2";
  const author = "ICheered";
  const authorLink = "https://icheered.nl";

  // PARAMETERS
  // let sourceFolderPath = "/home/tjbakker/Documents/dev/vscode/go/wails/testfiles/source";
  // let targetFolderPath = "/home/tjbakker/Documents/dev/vscode/go/wails/testfiles/watermarked";
  // let watermarkPath = "/home/tjbakker/Documents/dev/vscode/go/wails/testfiles/watermark.png";

  let sourceFolderPath = "";
  let targetFolderPath = "";
  let watermarkPath = "";

  let watermarkOpacity = 50;
  let watermarkPosition = "bottom-right";
  let watermarkScale = 20;

  // INFORMATION
  let changedSettings = true;

  let numberOfSourceFiles = 0;
  let numberOfTargetFiles = 0;

  let watermarkpreviewImage;
  let watermarkedpreviewImage;

  let returnval = "";

  // FUNCTIONS
  let mainbuttontext = "";
  $: {
    if (sourceFolderPath == "" || targetFolderPath == "" || watermarkPath == "") {
      mainbuttontext = "Select source, target and watermark";
    } else if (numberOfSourceFiles == 0) {
      mainbuttontext = "No images found in source folder!";
    } else if (changedSettings) {
      mainbuttontext = "Generate preview";
    } else {
      mainbuttontext = "Start processing!";
    }
  }

  function generatePreview() {
    console.log("Generating preview");
    FetchPreview(
      watermarkPath,
      sourceFolderPath,
      targetFolderPath,
      watermarkPosition,
      watermarkOpacity / 100,
      watermarkScale / 100
    ).then((result) => {
      if (result.status && result.status == "error") {
        alert(result.message);
      } else if (result.status && result.status == "success") {
        fetch(result.message)
          .then((response) => response.blob())
          .then((blob) => {
            const reader = new FileReader();
            reader.addEventListener("load", function () {
              watermarkedpreviewImage.setAttribute("src", reader.result);
            });
            reader.readAsDataURL(blob);
          });
      }
    });
  }

  function processFiles() {
    if (!sourceFolderPath || !targetFolderPath || !watermarkPath) {
      alert("Please select all the required fields");
      return;
    }
    ProcessImages(
      watermarkPath,
      sourceFolderPath,
      targetFolderPath,
      watermarkPosition,
      watermarkOpacity,
      watermarkScale
    ).then((result) => {
      if (result.status && result.status == "error") {
        alert(result.message);
      } else {
        alert(result.message);
      }
    });
  }

  function mainButtonFunction() {
    if (sourceFolderPath == "" || targetFolderPath == "" || watermarkPath == "") {
      alert("Please select all the required fields");
    } else if (numberOfSourceFiles == 0) {
      alert("No images found in source folder!");
    } else if (changedSettings) {
      generatePreview();
    } else {
      processFiles();
    }
  }
</script>

<main>
  <Header {title} {version} {author} {authorLink} />

  <div class="container">
    <div class="pathconfig">
      <Pathconfig
        bind:sourceFolderPath
        bind:targetFolderPath
        bind:watermarkPath
        bind:numberOfSourceFiles
        bind:numberOfTargetFiles
        bind:watermarkpreviewImage
      />
    </div>

    <div class="mainview col">
      <img bind:this={watermarkedpreviewImage} src="" alt="Preview of the result" />
      <button on:click={mainButtonFunction}>{mainbuttontext}</button>
    </div>

    <div class="settings col">
      <Settings bind:watermarkOpacity bind:watermarkScale bind:watermarkPosition />
    </div>
  </div>
</main>

<style>
  main {
    background: darkslateblue;
    height: 100%;
    width: 100%;
  }

  .container {
    display: flex;
    flex-direction: row;

    height: 100%;
    width: 100%;
  }

  .pathconfig {
    width: 25%;
  }
  .mainview {
    width: 50%;
  }

  .settings {
    width: 25%;
  }

  .col {
    display: flex;
    flex-direction: column;
  }
</style>
