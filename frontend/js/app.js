function interceptClickEvent(e) {
    if (e.target.tagName === 'A') {
        e.preventDefault();

        runtime.BrowserOpenURL(e.target.getAttribute('href'));
    }
}

//listen for link click events at the document level
if (document.addEventListener) {
    document.addEventListener('click', interceptClickEvent);
} else if (document.attachEvent) {
    document.attachEvent('onclick', interceptClickEvent);
}