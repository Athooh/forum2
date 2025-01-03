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
                            <div class="post">
                                <h4><a href="/view-post?id=${post.ID}">${post.Title}</a></h4>
                                <div class="post-image">
                                    <img src="${post.ImageURL ? post.ImageURL : '/static/images/default-post.jpg'}" alt="Post Image" class="post-image-preview">
                                </div>
                                <p>${post.Preview}</p>
                            </div>
                            <div class="reaction">
                                <i class="btn-like" data-id="${post.ID}"><i class="fa-regular fa-thumbs-up"></i> <span>${post.Likes}</span></i>
                                <i class="btn-dislike" data-id="${post.ID}"><i class="fa-regular fa-thumbs-down"></i> <span>${post.Dislikes}</span></i>
                                <i class="btn-comment" data-id="${post.ID}"><i class="fa-regular fa-message"></i> <span>${post.CommentsCount}</span></i>
                                <i class="fa-solid fa-share-nodes"></i>
                            </div>
                            ${post.IsOwner ? `
                                <div class="post-options">
                                    <button class="btn-edit" data-id="${post.ID}" style="border: none">Edit</button>
                                    <button class="btn-delete" data-id="${post.ID}" style="border: none">Delete</button>
                                </div>
                            ` : ''}
                            <div class="comment-form-container" id="comment-form-${post.ID}" style="display: none;">
                                <div class="guest-comment-prompt">
                                    <p>Click to view post and join the discussion</p>
                                </div>
                            </div>
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
    const editButtons = document.querySelectorAll('.btn-edit');
    const deleteButtons = document.querySelectorAll('.btn-delete');

    // Initialize like/dislike handlers
    likeButtons.forEach(button => {
        button.removeEventListener('click', handleLike);
        button.addEventListener('click', handleLike);
    });

    dislikeButtons.forEach(button => {
        button.removeEventListener('click', handleDislike);
        button.addEventListener('click', handleDislike);
    });

    // Initialize edit handlers
    editButtons.forEach(button => {
        button.removeEventListener('click', handleEdit);
        button.addEventListener('click', handleEdit);
    });

    // Initialize delete handlers
    deleteButtons.forEach(button => {
        button.removeEventListener('click', handleDelete);
        button.addEventListener('click', handleDelete);
    });

    // Initialize comment handlers
    commentButtons.forEach(button => {
        button.removeEventListener('click', handleComment);
        button.addEventListener('click', handleComment);
    });
}

// Handler functions
function handleLike(e) {
    e.preventDefault();
    const postId = this.getAttribute('data-id');
    fetch(`/like-post?id=${postId}`, { method: 'POST' })
        .then(response => response.json())
        .then(data => updateReactionCounts(postId, data.likes, data.dislikes))
        .catch(err => console.error('Error:', err));
}

function handleDislike(e) {
    e.preventDefault();
    const postId = this.getAttribute('data-id');
    fetch(`/dislike-post?id=${postId}`, { method: 'POST' })
        .then(response => response.json())
        .then(data => updateReactionCounts(postId, data.likes, data.dislikes))
        .catch(err => console.error('Error:', err));
}

function handleEdit(e) {
    e.preventDefault();
    const postId = this.getAttribute('data-id');
    window.location.href = `/edit-post?id=${postId}`;
}

function handleDelete(e) {
    e.preventDefault();
    const postId = this.getAttribute('data-id');
    if (confirm('Are you sure you want to delete this post?')) {
        fetch(`/delete-post?id=${postId}`, { method: 'POST' })
            .then(response => {
                if (response.ok) {
                    document.getElementById(`post-${postId}`).remove();
                }
            })
            .catch(err => console.error('Error:', err));
    }
}

function handleComment(e) {
    e.preventDefault();
    const postId = this.getAttribute('data-id');
    
    // Always redirect to view-post page when comment button is clicked
    window.location.href = `/view-post?id=${postId}#comments-section`;
} 