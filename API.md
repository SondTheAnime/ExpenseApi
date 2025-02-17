# API do Rastreador de Despesas

Crie uma API para um aplicativo de controle de despesas. Esta API deve permitir que os usuários criem, leiam, atualizem e excluam despesas. Os usuários devem poder se inscrever e fazer login no aplicativo. Cada usuário deve ter seu próprio conjunto de despesas.

![API.png](API.png)

## Características
Aqui estão os recursos que você deve implementar na sua API do Expense Tracker:

- Cadastre-se como um novo usuário.
- Gere e valide JWTs para manipular autenticação e sessão de usuário.
- Liste e filtre suas despesas passadas. Você pode adicionar os seguintes filtros:
  - Semana passada
  - Mês passado
  - Últimos 3 meses
  - Personalizado (para especificar uma data de início e término de sua escolha).
- Adicionar uma nova despesa
- Remover despesas existentes
- Atualizar despesas existentes

## Restrições
Você pode usar qualquer linguagem de programação e framework de sua escolha. Você pode usar um banco de dados de sua escolha para armazenar os dados. Você pode usar qualquer ORM ou biblioteca de banco de dados para interagir com o banco de dados.

Aqui estão algumas restrições que você deve seguir:

- Você usará JWT (JSON Web Token) para proteger os endpoints e identificar o solicitante.
- Para as diferentes categorias de despesas, você pode usar a seguinte lista (sinta-se à vontade para decidir como implementar isso como parte do seu modelo de dados):
  - Mantimentos
  - Lazer
  - Eletrônica  
  - Utilitários
  - Roupas
  - Saúde
  - Outros