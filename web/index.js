const BASE_URL = "/search?q=";
const DOCUMENT_INFO_URL = "/document/";

document.addEventListener("DOMContentLoaded",
 (event) => {
    const searchButtton = document.getElementById("search-btn");
    const resultsWrapper = document.getElementsByClassName("results-wrapper")[0];
    const searchQuery = document.getElementById("search-query")
    const resultsDiv = document.getElementById("search-results");
    const documentInfoDiv = document.getElementById("document-info-div");

    displayResults = (data) => {
        let resultsHtml = `<div class="summary">${data.result_count} results fetched in ${data.duration} secs.</div>`;
        resultsDiv.innerHTML = resultsHtml;
        data.results.forEach( result => {
            resultsHtml += `<div><div class="result-title"><a href="#document-info" data-id="${result.document_id}" class="results-link">${result.document}</a> <div class="summary">${result.ranking}</div></div></div>`;
        });
        resultsDiv.innerHTML = resultsHtml;
        resultsWrapper.style.visibility = "visible";

        let links = document.getElementsByClassName("results-link");
        for (let link of links) {
            link.addEventListener(
                'click',
                (evt) => {
                    documentInfoDiv.innerHTML = 'Loading...';
                    displayDocumentInfo(evt.srcElement.getAttribute('data-id'))
                }
            )
        }
    }

    getResults = async (query) => {
        const searchUrl = BASE_URL + encodeURIComponent(query)
        resultsWrapper.style.visibility = "hidden";
        const response = await fetch(searchUrl, {
            method: "GET"
        });
        return await response.json();
    }

    displayDocumentInfo = async (docId) => {
        const response = await fetch(`${DOCUMENT_INFO_URL}${docId}`, {
            method: "GET"
        });
        const data = await response.json();

        docInfoHtml = `<h3>Document ${data.id}</h3>
        <table>
            <tr><th>Id</th><td>${data.id}</td></tr>
            <tr><th>Path</th><td>${data.path}</td></tr>
            <tr><th>Parser</th><td>${data.parser}</td></tr>
            <tr><th>Terms</th><td>${data.total_terms}</td></tr>
        </table>
        <h3>Keywords</h3>
        <table><tr><th>Keyword</th><th>Count</th><th>Density</th></tr>`;

        data.keywords?.forEach( keyword => {
            docInfoHtml += `<tr><td>${keyword.keyword}</td><td>${keyword.keyword_count}</td><td>${keyword.keyword_density}</td></tr>`;
        });
        docInfoHtml += `</table>`;

        documentInfoDiv.innerHTML = docInfoHtml;
    }

    searchButtton.addEventListener("click", async () => {
        displayResults(await getResults(searchQuery.value))
    });
 }
)

console.log("Application has been initialized");