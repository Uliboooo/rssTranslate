document.addEventListener('DOMContentLoaded', () => {
    const API_URL = 'http://localhost:8080/api/rts';
    const contentsList = document.getElementById('contents-list');

    const fetchContents = async () => {
        try {
            const response = await fetch(API_URL);
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            const data = await response.json();
            return data.list;
        } catch (error) {
            console.error('Error fetching contents:', error);
            contentsList.innerHTML = '<li>Error fetching contents.</li>';
            return [];
        }
    };

    const renderContents = (contents) => {
        console.log(contents)
        if (contents.length === 0) {
            contentsList.innerHTML = '<li>No contents available.</li>';
            return;
        }

        const listItems = contents.map(item => {
            const li = document.createElement('li');
            const a = document.createElement('a');
            a.href = item.link;
            a.target = '_blank';
            a.textContent = item.title;
            const date = document.createElement('span');
            date.textContent = new Date(item.date).toLocaleString();
            li.appendChild(a);
            li.appendChild(document.createElement('br'));
            li.appendChild(date);
            return li;
        });

        contentsList.innerHTML = '';
        listItems.forEach(item => contentsList.appendChild(item));
    };

    const init = async () => {
        const contents = await fetchContents();
        renderContents(contents);
    };

    init();
});
