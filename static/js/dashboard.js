// inline editing for posts using javascript
document.addEventListener("DOMContentLoaded", () => {
    const editButtons = document.querySelectorAll(".edit-post-btn");

    editButtons.forEach((btn) => {
        btn.addEventListener("click", (event) => {
            const postId = btn.dataset.id;
            const postCard = document.getElementById(`post-${postId}`);
            const postTitle = postCard.querySelector("h4").textContent;
            const postContent = postCard.querySelector("p").textContent;

            // Replace content with an editable form
            postCard.innerHTML = `
                <form action="/edit-post?id=${postId}" method="POST" class="edit-post-form">
                    <input type="text" name="title" value="${postTitle}" required>
                    <textarea name="content" required>${postContent}</textarea>
                    <button type="submit">Save</button>
                    <button type="button" class="cancel-edit-btn">Cancel</button>
                </form>
            `;

            // Add cancel button functionality
            const cancelButton = postCard.querySelector(".cancel-edit-btn");
            cancelButton.addEventListener("click", () => {
                postCard.innerHTML = `
                    <div class="post">
                        <h4>${postTitle}</h4>
                        <p>${postContent}</p>
                        <div class="reaction">
                            <i class="fa-regular fa-thumbs-up"> 42</i>
                            <i class="fa-regular fa-thumbs-down"> 2</i>
                            <i class="fa-regular fa-message"> 12</i>
                            <i class="fa-solid fa-share-nodes"></i>
                        </div>
                    </div>
                `;
            });
        });
    });
});
