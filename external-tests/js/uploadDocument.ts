/*
This test runs through the upload document workflow in typescript hitting endpoints
to the server running on localhost.
*/

const BASE_URL = "http://localhost:8000"
const CUSTOMER_ID = 7
const PARENT_FOLDER_ID = 14

async function uploadDocument(filename: string) {

    // Step 1: Read the file (adapted for Bun)
    const file = Bun.file(`../../resources/${filename}`);
    const contents = await file.arrayBuffer();

    // Create a SHA-256 hash of the file content
    const hasher = new Bun.CryptoHasher("sha256");
    hasher.update(contents);
    const signature = hasher.digest("base64")

    // create a json body for the request to process
    const fileUploadbody = {
        filename: filename,
        mime: file.type,
        signature: signature,
        size: file.size,
        parentId: PARENT_FOLDER_ID,
    }

    // Step 2: Get pre-signed URL
    const presignedUrlResponse = await fetch(`${BASE_URL}/customers/${CUSTOMER_ID}/generatePresignedUrl`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(fileUploadbody),
    });

    if (!presignedUrlResponse.ok) {
        throw new Error('Failed to get pre-signed URL');
    }

    const uploadData = await presignedUrlResponse.json();

    // Step 3: Upload the document
    const uploadResponse = await fetch(uploadData.uploadUrl, {
        method: uploadData.method,
        body: contents,
        headers: {
            'Content-Type': file.type,
        },
    });

    if (!uploadResponse.ok) {
        console.log(uploadResponse)
        const d = await uploadResponse.text()
        console.log(d)
        throw new Error('Failed to upload document');
    }

    // Step 4: Notify successful upload
    const notifyResponse = await fetch(`${BASE_URL}/customers/${CUSTOMER_ID}/documents/${uploadData.documentId}/validate`, {
        method: 'PUT',
    });

    if (!notifyResponse.ok) {
        throw new Error('Failed to notify successful upload');
    }

    console.log('Document uploaded and notification sent successfully');
}

uploadDocument('../resources/file1.txt').catch(console.error);
