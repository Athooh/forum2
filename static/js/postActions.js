// likes and  dislike buttons functionalities
document.addEventListener('DOMContentLoaded', () => {
    const likeButtons = document.querySelectorAll('.btn-like');
    const dislikeButtons = document.querySelectorAll('.btn-dislike');

    // Handle Like Button Click
    likeButtons.forEach(button => {
        button.addEventListener('click', () => {
            const postId = button.getAttribute('data-id');
            fetch(`/like-post?id=${postId}`, { method: 'POST' })
                .then(response => response.json())
                .then(data => {
                    button.querySelector('span').innerText = data.likes;
                    const dislikeButton = button.nextElementSibling;
                    dislikeButton.querySelector('span').innerText = data.dislikes;
                })
                .catch(err => console.error('Error:', err));
        });
    });

    // Handle Dislike Button Click
    dislikeButtons.forEach(button => {
        button.addEventListener('click', () => {
            const postId = button.getAttribute('data-id');
            fetch(`/dislike-post?id=${postId}`, { method: 'POST' })
                .then(response => response.json())
                .then(data => {
                    button.querySelector('span').innerText = data.dislikes;
                    const likeButton = button.previousElementSibling;
                    likeButton.querySelector('span').innerText = data.likes;
                })
                .catch(err => console.error('Error:', err));
        });
    });
});
