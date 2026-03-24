describe('Login Page', () => {
    beforeEach(() => {
        cy.visit('/login');
    });

    it('should display the login form', () => {
        cy.get('input[type="email"], input[name="email"], input[placeholder*="email" i]').should('be.visible');
        cy.get('input[type="password"], input[name="password"], input[placeholder*="password" i]').should('be.visible');
    });

    it('should fill in email and password fields', () => {
        cy.get('input[type="email"], input[name="email"], input[placeholder*="email" i]')
            .type('testuser@ufl.edu')
            .should('have.value', 'testuser@ufl.edu');

        cy.get('input[type="password"], input[name="password"], input[placeholder*="password" i]')
            .type('TestPassword123')
            .should('have.value', 'TestPassword123');
    });

    it('should click the login/submit button', () => {
        cy.get('input[type="email"], input[name="email"], input[placeholder*="email" i]')
            .type('testuser@ufl.edu');

        cy.get('input[type="password"], input[name="password"], input[placeholder*="password" i]')
            .type('TestPassword123');

        cy.get('button[type="submit"], button').contains(/sign in|log in|login|submit/i).click();
    });

    it('should show a link to the register page', () => {
        cy.get('a[href*="register"], a[routerLink*="register"]').should('exist');
    });
});
