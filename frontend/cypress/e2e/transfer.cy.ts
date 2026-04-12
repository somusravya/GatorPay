describe('Transfer Page', () => {
    beforeEach(() => {
        // Mock authentication state
        window.localStorage.setItem('gatorpay_token', 'fake-jwt-token');
        
        // Mock user profile response
        cy.intercept('GET', '**/api/v1/auth/me', {
            statusCode: 200,
            body: {
                success: true,
                data: {
                    user: { id: 'user-1', first_name: 'Test', last_name: 'User', email: 'test@example.com' },
                    wallet: { balance: 1540.50, currency: 'USD' }
                }
            }
        }).as('getMe');

        // Mock recent contacts
        cy.intercept('GET', '**/api/v1/transfer/contacts', {
            statusCode: 200,
            body: {
                success: true,
                data: [
                    { id: 'user-2', username: 'alice99', first_name: 'Alice', last_name: 'Smith', email: 'alice@example.com' }
                ]
            }
        }).as('getContacts');

        cy.visit('/transfer');
    });

    it('should display the transfer form', () => {
        cy.wait('@getMe');
        cy.get('input[formControlName="recipient"], input[name="recipient"], input[placeholder*="username" i]').should('be.visible');
        cy.get('input[formControlName="amount"], input[name="amount"], input[placeholder*="amount" i]').should('be.visible');
        cy.get('input[formControlName="note"], input[name="note"], input[placeholder*="note" i]').should('exist');
        cy.get('button[type="submit"], button').contains(/send|transfer/i).should('be.visible');
    });

    it('should show the user balance', () => {
        cy.wait('@getMe');
        cy.contains(/balance/i).should('be.visible');
        cy.contains(/\$1,540\.50|1540\.50/).should('be.visible');
    });

    it('should display recent contacts', () => {
        cy.wait('@getContacts');
        cy.contains(/alice99|Alice Smith/i).should('be.visible');
    });
});
