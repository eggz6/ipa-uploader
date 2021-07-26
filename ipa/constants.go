package ipa

const manifest = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>items</key>
	<array>
		<dict>
			<key>assets</key>
			<array>
				<dict>
					<key>kind</key>
					<string>software-package</string>
					<key>url</key>
					<string>${IPA_URL}</string>
				</dict>
				<dict>
					<key>kind</key>
					<string>display-image</string>
					<key>needs-shine</key>
					<integer>0</integer>
					<key>url</key>
					<string><![CDATA[${ICON_URL}]]></string>
					</dict>
				<dict>
					<key>kind</key>
					<string>full-size-image</string>
					<key>needs-shine</key>
					<true/>
					<key>url</key>
					<string><![CDATA[${ICON_URL}]]></string>
				</dict>
			</array>
			<key>metadata</key>
			<dict>
				<key>bundle-identifier</key>
				<string>${BUNDLE_ID}</string>
				<key>bundle-version</key>
				<string>${BUNDLE_VERSION}</string>
				<key>kind</key>
				<string>software</string>
				<key>title</key>
				<string>${TITLE}</string>
			</dict>
		</dict>
	</array>
</dict>
</plist>
`
const html = `
<!DOCTYPE html>
<html lang="en">
    <head>
        <meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
        <title>${TITLE}</title>
        <style>
            h1 {
                text-align: center;
		vertical-align: middle;
            }
            p {
                text-align: center;
            }
	        div {
		          position: absolute;
                  top: 50%;
                  left: 50%;
		          height: 30%;
                  width: 50%;
		          margin: -15% 0 0 -25%;
	        }

	    .round-button {
                display:block;
                height:80px;
                width:360px;
                line-height:80px;
                border: 2px solid #f5f5f5;
                border-radius: 10px;
                color:#f5f5f5;
                text-align:center;
                text-decoration:none;
                background: #464646;
                box-shadow: 0 0 3px gray;
                font-size:40px;
                font-weight:bold;
                margin: auto;
            }
            .round-button:hover {
                background: #262626;
            }
        </style>
    </head>
    <body>
	<div>
	        <img src="${ICON_URL}">
            <h1>${MESSAGE}</h1>
            <p><a href="itms-services://?action=download-manifest&url=${MANIFEST_URL}" class="round-button">Install</a></p>
	</div>
    </body>
</html>
`
