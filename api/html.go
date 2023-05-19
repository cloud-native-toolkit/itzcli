package api

const succesHTML = `<html>

<head>
    <meta charset="UTF-8" />
    <link rel="stylesheet" href="https://1.www.s81c.com/common/carbon-for-ibm-dotcom/tag/v1/latest/plex.css" />
    <link rel="stylesheet" href="https://unpkg.com/carbon-components/css/carbon-components.min.css">
    <style>
        /* Suppress custom element until styles are loaded */
        bx-header:not(:defined) {
            display: none;
        }

        #app {
            font-family: 'IBM Plex Sans', 'Helvetica Neue', Arial, sans-serif;
            width: 100%;
            margin-top: 2rem;
        }
    </style>
    <script type="module" src="https://1.www.s81c.com/common/carbon/web-components/tag/latest/ui-shell.min.js"></script>
    <script type="module" src="	https://unpkg.com/carbon-components/scripts/carbon-components.min.js"></script>
    <script type="module"
        src="https://1.www.s81c.com/common/carbon/web-components/tag/latest/notification.min.js"></script>
</head>

<body>
    <div id="app">
        <bx-header aria-label="IBM Platform Name">
            <bx-header-menu-button button-label-active="Close menu"
                button-label-inactive="Open menu"></bx-header-menu-button>
            <bx-header-name href="javascript:void 0" prefix="IBM">Technology Zone</bx-header-name>
        </bx-header>
        <main class="bx--content bx-ce-demo-devenv--ui-shell-content">
            <bx-inline-notification style="min-width: 30rem; margin-bottom: .5rem" hide-close-button kind="success"
                title="ITZCLI Login Success!"
                subtitle="Hello {{.Preferredfirstname}} {{.Preferredlastname}}! You are now authenticated with ITZCLI. You may close this browser and return to the terminal.">
            </bx-inline-notification>
        </main>
    </div>
</body>

</html>`


const errorHTML = `<html>

<head>
    <meta charset="UTF-8" />
    <link rel="stylesheet" href="https://1.www.s81c.com/common/carbon-for-ibm-dotcom/tag/v1/latest/plex.css" />
    <link rel="stylesheet" href="https://unpkg.com/carbon-components/css/carbon-components.min.css">
    <style>
        /* Suppress custom element until styles are loaded */
        bx-header:not(:defined) {
            display: none;
        }

        #app {
            font-family: 'IBM Plex Sans', 'Helvetica Neue', Arial, sans-serif;
            width: 100%;
            margin-top: 2rem;
        }
    </style>
    <script type="module" src="https://1.www.s81c.com/common/carbon/web-components/tag/latest/ui-shell.min.js"></script>
    <script type="module" src="	https://unpkg.com/carbon-components/scripts/carbon-components.min.js"></script>
    <script type="module"
        src="https://1.www.s81c.com/common/carbon/web-components/tag/latest/notification.min.js"></script>
</head>

<body>
    <div id="app">
        <bx-header aria-label="IBM Platform Name">
            <bx-header-menu-button button-label-active="Close menu"
                button-label-inactive="Open menu"></bx-header-menu-button>
            <bx-header-name href="javascript:void 0" prefix="IBM">Technology Zone</bx-header-name>
        </bx-header>
        <main class="bx--content bx-ce-demo-devenv--ui-shell-content">
            <bx-inline-notification style="min-width: 30rem; margin-bottom: .5rem" hide-close-button kind="error"
                title="ITZCLI Login ERROR!"
                subtitle="Error trying to login via the CLI. Please close the browser and try again from the terminal.">
            </bx-inline-notification>
        </main>
    </div>
</body>

</html>`