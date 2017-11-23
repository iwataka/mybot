function deleteRow(button, modal) {
    var tr = button.parentNode.parentNode;
    if (tr.className === 'deleted') {
        for (var i = 0; i < tr.childNodes.length; i++) {
            var child = tr.childNodes[i]
            if (child.type === 'hidden' && child.nodeName === 'INPUT') {
                child.value = 'false';
                tr.className = '';
                button.innerHTML = 'Delete';
            }
        }
    } else {
        for (var i = 0; i < tr.childNodes.length; i++) {
            var child = tr.childNodes[i]
            if (child.type === 'hidden' && child.nodeName === 'INPUT') {
                child.value = 'true';
                tr.className = 'deleted';
                button.innerHTML = 'Revert';
            }
        }
    }
}

