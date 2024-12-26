document.addEventListener('DOMContentLoaded', () => {
    const categoryLinks = document.querySelectorAll('.category-link');
    const postsContainer = document.querySelector('.posts');

    categoryLinks.forEach(link => {
        link.addEventListener('click', async (e) => {
            e.preventDefault();
            console.log('Category link clicked'); // Debug log
            
            const category = link.getAttribute('data-category');
            console.log('Selected category:', category); // Debug log
            
            try {
                const response = await fetch(`/posts-by-category?category=${encodeURIComponent(category)}`);
                console.log('Response status:', response.status); // Debug log
                
                if (!response.ok) {
                    throw new Error(`HTTP error! status: ${response.status}`);
                }
                
                const data = await response.json();
                console.log('Received data:', data); // Debug log

                if (data.posts && Array.isArray(data.posts)) {
                    // Update posts container
                    const postsHTML = data.posts.map(post => `
                        <div class="post-card" id="post-${post.ID}">
                            <div class="profile">
                                <div class="avatar">
                                    <i class="fa-regular fa-user"></i>
                                    <div class="avatar-name">
                                        <p>${post.Username}</p>
                                        <p>${post.CreatedAtHuman}</p>
                                    </div>
                                </div>
                                <div class="category-tag">
                                    <p>${post.Category}</p>
                                </div>
                            </div>
                            <a href="/view-post?id=${post.ID}">
                            <div class="post">
                                <h4>${post.Title}</h4>
                                <div class="post-image">
                                    <img src="${post.ImageURL || '/static/images/default-post.jpg'}" alt="Post Image" class="post-image-preview">
                                </div>
                                <p>${post.Preview}</p>
                                <div class="reaction">
                                    <i class="btn-like" data-id="${post.ID}"><i class="fa-regular fa-thumbs-up"></i> <span>${post.Likes}</span></i>
                                    <i class="btn-dislike" data-id="${post.ID}"><i class="fa-regular fa-thumbs-down"></i> <span>${post.Dislikes}</span></i>
                                    <i class="btn-comment" data-id="${post.ID}"><i class="fa-regular fa-message"></i> <span>${post.CommentsCount}</span></i>
                                    <i class="fa-solid fa-share-nodes"></i>
                                </div>
                            </div>
                            </a>
                        </div>
                    `).join('');

                    // Update the posts container with new content
                    postsContainer.innerHTML = `
                        <h3>${category} discussions</h3>
                        <p>Join the conversation and share your thoughts</p>
                        <div class="create-post">
                            <a href="/create-post" class="btn-create-post">Create New Post</a>
                        </div>
                        ${postsHTML}
                    `;

                    // Reinitialize event listeners
                    initializeEventListeners();
                }
            } catch (error) {
                console.error('Error fetching posts:', error);
            }
        });
    });
});

// Helper function to reinitialize event listeners
function initializeEventListeners() {
    // Reinitialize likes/dislikes
    const likeButtons = document.querySelectorAll('.btn-like');
    const dislikeButtons = document.querySelectorAll('.btn-dislike');
    const commentButtons = document.querySelectorAll('.btn-comment');

    likeButtons.forEach(button => {
        button.addEventListener('click', (e) => {
            e.preventDefault();
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

    dislikeButtons.forEach(button => {
        button.addEventListener('click', (e) => {
            e.preventDefault();
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

    commentButtons.forEach(button => {
        button.addEventListener('click', () => {
            const postId = button.getAttribute('data-id');
            const commentForm = document.getElementById(`comment-form-${postId}`);
            if (commentForm) {
                commentForm.style.display = commentForm.style.display === 'none' ? 'block' : 'none';
            }
        });
    });
} 