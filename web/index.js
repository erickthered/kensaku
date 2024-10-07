const BASE_URL = "/search?q=";

document.addEventListener("DOMContentLoaded", 
 (event) => {
    const searchButtton = document.getElementById("search-btn");

    searchButtton.addEventListener("click", async () => {
        const resultsWrapper = document.getElementsByClassName("results-wrapper")[0];
        const searchQuery = document.getElementById("search-query")
        const searchUrl = BASE_URL + encodeURIComponent(searchQuery.value)
        resultsWrapper.style.visibility = "hidden";

        const response = await fetch(searchUrl, {
            method: "GET"
        })
        const data= await response.json();

        const resultsDiv = document.getElementById("search-results");
        let resultsHtml = `<div class="summary">${data.result_count} results fetched in ${data.duration} secs.</div>`;
        resultsDiv.innerHTML = resultsHtml;

        data.results.forEach( result => {
            resultsHtml += `<div><div class="result-title"><a href="file:///${result.document}">${result.document}</a> <div class="summary">${result.ranking}</div></div></div>`;
        });
        resultsDiv.innerHTML = resultsHtml;
        resultsWrapper.style.visibility = "visible";
    });
 }
)

console.log("Application has been initialized");